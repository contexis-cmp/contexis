"""
Local AI Provider for Contexis
Uses Phi-3.5-Mini for local text generation without external API calls.
"""

import json
import os
import sys

import logging
from typing import Dict, List, Optional, Any
from pathlib import Path

try:
    from transformers import AutoTokenizer, AutoModelForCausalLM
    import torch
    from sentence_transformers import SentenceTransformer
except ImportError as e:
    raise ImportError(f"Local provider requires transformers and torch: {e}")

logger = logging.getLogger(__name__)


class LocalAIProvider:
    """Local AI provider using Phi-3.5-Mini for text generation."""
    
    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.model_name = config.get("model", "microsoft/DialoGPT-medium")
        self.temperature = config.get("temperature", 0.1)
        self.max_tokens = config.get("max_tokens", 1000)
        self.device = config.get("device", "auto")
        self.load_in_8bit = config.get("load_in_8bit", False)
        self.load_in_4bit = config.get("load_in_4bit", False)
        
        # Model cache directory
        cache_dir = config.get("model_cache", {}).get("directory", "./data/models")
        self.cache_dir = Path(cache_dir)
        self.cache_dir.mkdir(parents=True, exist_ok=True)
        
        self._model = None
        self._tokenizer = None
        
    def _load_model(self):
        """Load the Phi-3.5-Mini model."""
        if self._model is not None and self._tokenizer is not None:
            return
            
        logger.info(f"Loading local model: {self.model_name}")
        
        # Determine device
        if self.device == "auto":
            device = "cuda" if torch.cuda.is_available() else "cpu"
        else:
            device = self.device
            
        logger.info(f"Using device: {device}")
        
        # Load quantization config if needed
        quantization_config = None
        if self.load_in_8bit or self.load_in_4bit:
            try:
                from transformers import BitsAndBytesConfig
                quantization_config = BitsAndBytesConfig(
                    load_in_8bit=self.load_in_8bit,
                    load_in_4bit=self.load_in_4bit,
                )
            except ImportError:
                logger.warning("BitsAndBytes not available, loading model in full precision")
        
        # Load model and tokenizer
        try:
            self._tokenizer = AutoTokenizer.from_pretrained(
                self.model_name,
                cache_dir=self.cache_dir,
                trust_remote_code=True
            )
            
            self._model = AutoModelForCausalLM.from_pretrained(
                self.model_name,
                cache_dir=self.cache_dir,
                torch_dtype=torch.float16 if device == "cuda" else torch.float32,
                quantization_config=quantization_config,
                trust_remote_code=True,
                device_map="auto" if device == "cuda" else None,
                attn_implementation="eager",
            )
            
            if device == "cpu":
                self._model = self._model.to(device)
                
            logger.info("Local model loaded successfully")
            
        except Exception as e:
            logger.error(f"Failed to load local model: {e}")
            raise
    
    def generate(self, prompt: str, **kwargs) -> str:
        """Generate text using the local model."""
        self._load_model()
        
        # Merge config with kwargs
        generation_config = {
            "temperature": self.temperature,
            "max_new_tokens": self.max_tokens,
            "do_sample": True,
            "pad_token_id": self._tokenizer.eos_token_id,
            "use_cache": False,
        }
        generation_config.update(kwargs)
        
        try:
            # Format prompt for Phi-3.5-Mini
            formatted_prompt = self._format_prompt(prompt)
            
            # Tokenize input
            inputs = self._tokenizer(formatted_prompt, return_tensors="pt")
            
            # Move to device
            if hasattr(self._model, 'device'):
                inputs = {k: v.to(self._model.device) for k, v in inputs.items()}
            
            # Generate using model directly (simpler approach to avoid DynamicCache issues)
            with torch.no_grad():
                outputs = self._model.generate(
                    **inputs,
                    **generation_config,
                )
            
            # Decode the generated tokens
            generated_text = self._tokenizer.decode(outputs[0], skip_special_tokens=True)
            
            # Remove the input prompt from the output
            if generated_text.startswith(formatted_prompt):
                generated_text = generated_text[len(formatted_prompt):]
                
            return generated_text.strip()
            
        except Exception as e:
            logger.error(f"Generation failed: {e}")
            raise
    
    def _format_prompt(self, prompt: str) -> str:
        """Format prompt for the model."""
        # Simple prompt format for DialoGPT
        return prompt
    
    def chat(self, messages: List[Dict[str, str]], **kwargs) -> str:
        """Generate chat response from message history."""
        # Convert messages to a single prompt
        prompt = self._messages_to_prompt(messages)
        return self.generate(prompt, **kwargs)
    
    def _messages_to_prompt(self, messages: List[Dict[str, str]]) -> str:
        """Convert message history to a single prompt."""
        formatted_messages = []
        
        for message in messages:
            role = message.get("role", "user")
            content = message.get("content", "")
            
            if role == "user":
                formatted_messages.append(f"<|user|>\n{content}<|end|>")
            elif role == "assistant":
                formatted_messages.append(f"<|assistant|>\n{content}<|end|>")
            elif role == "system":
                # System messages can be prepended as user messages
                formatted_messages.append(f"<|user|>\nSystem: {content}<|end|>")
        
        # Add the assistant prefix for the response
        formatted_messages.append("<|assistant|>\n")
        
        return "\n".join(formatted_messages)


class LocalEmbeddingsProvider:
    """Local embeddings provider using Sentence Transformers."""
    
    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.model_name = config.get("model", "all-MiniLM-L6-v2")
        self.device = config.get("device", "auto")
        
        # Model cache directory
        cache_dir = config.get("model_cache", {}).get("directory", "./data/models")
        self.cache_dir = Path(cache_dir)
        self.cache_dir.mkdir(parents=True, exist_ok=True)
        
        self._model = None
        
    def _load_model(self):
        """Load the sentence transformer model."""
        if self._model is not None:
            return
            
        logger.info(f"Loading local embeddings model: {self.model_name}")
        
        # Determine device
        if self.device == "auto":
            device = "cuda" if torch.cuda.is_available() else "cpu"
        else:
            device = self.device
            
        logger.info(f"Using device: {device}")
        
        try:
            self._model = SentenceTransformer(
                self.model_name,
                cache_folder=self.cache_dir
            )
            self._model = self._model.to(device)
            logger.info("Local embeddings model loaded successfully")
            
        except Exception as e:
            logger.error(f"Failed to load local embeddings model: {e}")
            raise
    
    def embed(self, texts: List[str]) -> List[List[float]]:
        """Generate embeddings for a list of texts."""
        self._load_model()
        
        try:
            embeddings = self._model.encode(texts, convert_to_tensor=False)
            return embeddings.tolist()
        except Exception as e:
            logger.error(f"Embedding generation failed: {e}")
            raise
    
    def embed_single(self, text: str) -> List[float]:
        """Generate embedding for a single text."""
        return self.embed([text])[0]


# Factory functions for easy integration
def create_local_ai_provider(config: Dict[str, Any]) -> LocalAIProvider:
    """Create a local AI provider instance."""
    return LocalAIProvider(config)


def create_local_embeddings_provider(config: Dict[str, Any]) -> LocalEmbeddingsProvider:
    """Create a local embeddings provider instance."""
    return LocalEmbeddingsProvider(config)


# Minimal CLI runner: read JSON from stdin and write JSON to stdout
# Input: {"prompt": "...", "params": {"MaxNewTokens": 256, ...}}
# Output: {"output": "..."}

def _main():
    try:
        raw = sys.stdin.read()
        data = json.loads(raw) if raw else {}
        prompt = data.get("prompt", "")
        params = data.get("params", {})
        max_new_tokens = int(params.get("MaxNewTokens", 256))

        config: Dict[str, Any] = {
            "model": os.getenv("CMP_LOCAL_MODEL_ID", "microsoft/Phi-3-mini-4k-instruct"),
            "device": os.getenv("CMP_LOCAL_DEVICE", "auto"),
            "load_in_8bit": os.getenv("CMP_LOCAL_LOAD_8BIT", "false").lower() == "true",
            "load_in_4bit": os.getenv("CMP_LOCAL_LOAD_4BIT", "false").lower() == "true",
            "temperature": float(os.getenv("CMP_LOCAL_TEMPERATURE", "0.1")),
            "max_tokens": int(os.getenv("CMP_LOCAL_MAX_TOKENS", str(max_new_tokens))),
            "model_cache": {"directory": os.getenv("CMP_MODEL_CACHE_DIR", "./data/models")},
        }
        provider = LocalAIProvider(config)
        output = provider.generate(prompt, max_new_tokens=max_new_tokens)
        print(json.dumps({"output": output}))
    except Exception as e:
        print(json.dumps({"error": str(e)}))
        sys.exit(1)


if __name__ == "__main__":
    _main()
