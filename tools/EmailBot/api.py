#!/usr/bin/env python3
"""
API Tool for CMP Agents
Provides HTTP API call capabilities with authentication and error handling
"""

import requests
import json
import logging
from typing import Dict, Any, Optional, List
from dataclasses import dataclass
from urllib.parse import urljoin
import time

logger = logging.getLogger(__name__)

@dataclass
class APIResponse:
    """Represents an API response"""
    status_code: int
    data: Dict[str, Any]
    headers: Dict[str, str]
    response_time: float
    success: bool

class APITool:
    """API tool for agents with authentication and rate limiting"""
    
    def __init__(self, base_url: str = "", api_key: Optional[str] = None):
        self.base_url = base_url
        self.api_key = api_key
        self.session = requests.Session()
        self.rate_limit_delay = 1.0  # seconds between requests
        
        # Set default headers
        self.session.headers.update({
            'User-Agent': 'CMP-Agent/1.0',
            'Content-Type': 'application/json',
        })
        
        if api_key:
            self.session.headers.update({
                'Authorization': f'Bearer {api_key}'
            })
    
    def _make_request(self, method: str, endpoint: str, data: Optional[Dict] = None, 
                     params: Optional[Dict] = None) -> APIResponse:
        """
        Make HTTP request with error handling and rate limiting
        
        Args:
            method: HTTP method (GET, POST, PUT, DELETE)
            endpoint: API endpoint
            data: Request body data
            params: Query parameters
            
        Returns:
            APIResponse object
        """
        start_time = time.time()
        
        try:
            # Rate limiting
            time.sleep(self.rate_limit_delay)
            
            url = urljoin(self.base_url, endpoint)
            
            logger.info(f"Making {method} request to {url}")
            
            response = self.session.request(
                method=method,
                url=url,
                json=data,
                params=params,
                timeout=30
            )
            
            response_time = time.time() - start_time
            
            # Parse response
            try:
                response_data = response.json() if response.content else {}
            except json.JSONDecodeError:
                response_data = {"text": response.text}
            
            api_response = APIResponse(
                status_code=response.status_code,
                data=response_data,
                headers=dict(response.headers),
                response_time=response_time,
                success=response.status_code < 400
            )
            
            logger.info(f"API request completed: {response.status_code} in {response_time:.3f}s")
            return api_response
            
        except requests.RequestException as e:
            logger.error(f"API request failed: {e}")
            return APIResponse(
                status_code=0,
                data={"error": str(e)},
                headers={},
                response_time=time.time() - start_time,
                success=False
            )
    
    def get(self, endpoint: str, params: Optional[Dict] = None) -> APIResponse:
        """Make GET request"""
        return self._make_request("GET", endpoint, params=params)
    
    def post(self, endpoint: str, data: Optional[Dict] = None) -> APIResponse:
        """Make POST request"""
        return self._make_request("POST", endpoint, data=data)
    
    def put(self, endpoint: str, data: Optional[Dict] = None) -> APIResponse:
        """Make PUT request"""
        return self._make_request("PUT", endpoint, data=data)
    
    def delete(self, endpoint: str) -> APIResponse:
        """Make DELETE request"""
        return self._make_request("DELETE", endpoint)
    
    def call_weather_api(self, city: str, api_key: Optional[str] = None) -> Dict[str, Any]:
        """
        Get weather information for a city
        
        Args:
            city: City name
            api_key: Weather API key (optional)
            
        Returns:
            Weather data dictionary
        """
        try:
            # Use OpenWeatherMap API as example
            weather_key = api_key or self.api_key
            if not weather_key:
                return {"error": "No API key provided for weather service"}
            
            params = {
                'q': city,
                'appid': weather_key,
                'units': 'metric'
            }
            
            response = self.get("https://api.openweathermap.org/data/2.5/weather", params)
            
            if response.success:
                return {
                    "city": city,
                    "temperature": response.data.get("main", {}).get("temp"),
                    "description": response.data.get("weather", [{}])[0].get("description"),
                    "humidity": response.data.get("main", {}).get("humidity"),
                    "wind_speed": response.data.get("wind", {}).get("speed")
                }
            else:
                return {"error": f"Weather API error: {response.status_code}"}
                
        except Exception as e:
            logger.error(f"Weather API call failed: {e}")
            return {"error": str(e)}
    
    def call_news_api(self, query: str, api_key: Optional[str] = None) -> List[Dict[str, Any]]:
        """
        Get news articles for a query
        
        Args:
            query: Search query
            api_key: News API key (optional)
            
        Returns:
            List of news articles
        """
        try:
            # Use NewsAPI as example
            news_key = api_key or self.api_key
            if not news_key:
                return [{"error": "No API key provided for news service"}]
            
            params = {
                'q': query,
                'apiKey': news_key,
                'language': 'en',
                'sortBy': 'publishedAt'
            }
            
            response = self.get("https://newsapi.org/v2/everything", params)
            
            if response.success:
                articles = response.data.get("articles", [])
                return [
                    {
                        "title": article.get("title"),
                        "description": article.get("description"),
                        "url": article.get("url"),
                        "published_at": article.get("publishedAt"),
                        "source": article.get("source", {}).get("name")
                    }
                    for article in articles[:5]  # Limit to 5 articles
                ]
            else:
                return [{"error": f"News API error: {response.status_code}"}]
                
        except Exception as e:
            logger.error(f"News API call failed: {e}")
            return [{"error": str(e)}]
    
    def call_translation_api(self, text: str, target_language: str, 
                           api_key: Optional[str] = None) -> Dict[str, Any]:
        """
        Translate text to target language
        
        Args:
            text: Text to translate
            target_language: Target language code (e.g., 'es', 'fr')
            api_key: Translation API key (optional)
            
        Returns:
            Translation result
        """
        try:
            # Use Google Translate API as example
            translate_key = api_key or self.api_key
            if not translate_key:
                return {"error": "No API key provided for translation service"}
            
            data = {
                'q': text,
                'target': target_language,
                'source': 'auto'
            }
            
            response = self.post(
                f"https://translation.googleapis.com/language/translate/v2?key={translate_key}",
                data
            )
            
            if response.success:
                translation = response.data.get("data", {}).get("translations", [{}])[0]
                return {
                    "original_text": text,
                    "translated_text": translation.get("translatedText"),
                    "target_language": target_language,
                    "detected_language": translation.get("detectedSourceLanguage")
                }
            else:
                return {"error": f"Translation API error: {response.status_code}"}
                
        except Exception as e:
            logger.error(f"Translation API call failed: {e}")
            return {"error": str(e)}
    
    def call_currency_api(self, from_currency: str, to_currency: str, 
                         amount: float = 1.0, api_key: Optional[str] = None) -> Dict[str, Any]:
        """
        Convert currency using exchange rate API
        
        Args:
            from_currency: Source currency code (e.g., 'USD')
            to_currency: Target currency code (e.g., 'EUR')
            amount: Amount to convert
            api_key: Currency API key (optional)
            
        Returns:
            Currency conversion result
        """
        try:
            # Use ExchangeRate-API as example
            currency_key = api_key or self.api_key
            if not currency_key:
                return {"error": "No API key provided for currency service"}
            
            response = self.get(
                f"https://v6.exchangerate-api.com/v6/{currency_key}/pair/{from_currency}/{to_currency}"
            )
            
            if response.success:
                rate = response.data.get("conversion_rate", 0)
                return {
                    "from_currency": from_currency,
                    "to_currency": to_currency,
                    "amount": amount,
                    "converted_amount": amount * rate,
                    "exchange_rate": rate,
                    "last_updated": response.data.get("time_last_update_utc")
                }
            else:
                return {"error": f"Currency API error: {response.status_code}"}
                
        except Exception as e:
            logger.error(f"Currency API call failed: {e}")
            return {"error": str(e)}

def main():
    """Test the API tool"""
    # Initialize with optional API key
    api_tool = APITool(api_key="your_api_key_here")
    
    # Test basic GET request
    response = api_tool.get("https://jsonplaceholder.typicode.com/posts/1")
    print(f"GET response: {response.status_code}")
    print(f"Data: {response.data}")
    
    # Test POST request
    post_data = {"title": "Test Post", "body": "Test content", "userId": 1}
    response = api_tool.post("https://jsonplaceholder.typicode.com/posts", post_data)
    print(f"\nPOST response: {response.status_code}")
    print(f"Data: {response.data}")
    
    # Test weather API (requires API key)
    # weather = api_tool.call_weather_api("London")
    # print(f"\nWeather: {weather}")
    
    # Test news API (requires API key)
    # news = api_tool.call_news_api("artificial intelligence")
    # print(f"\nNews: {news}")

if __name__ == "__main__":
    main()
