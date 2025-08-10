#!/usr/bin/env python3
"""
File System Tool for CMP Agents
Provides secure file read/write capabilities with access control
"""

import os
import json
import logging
from typing import Dict, Any, Optional, List, Union
from dataclasses import dataclass
from pathlib import Path
import hashlib
import mimetypes
from datetime import datetime
import shutil

logger = logging.getLogger(__name__)

@dataclass
class FileInfo:
    """Represents file information"""
    path: str
    name: str
    size: int
    modified: datetime
    mime_type: str
    is_directory: bool
    permissions: str

@dataclass
class FileOperation:
    """Represents a file operation result"""
    success: bool
    operation: str
    file_path: str
    data: Optional[Any] = None
    error: Optional[str] = None
    file_hash: Optional[str] = None

class FileSystemTool:
    """File system tool for agents with security and access control"""
    
    def __init__(self, base_path: str = ".", allowed_extensions: Optional[List[str]] = None):
        self.base_path = Path(base_path).resolve()
        self.allowed_extensions = allowed_extensions or ['.txt', '.md', '.json', '.yaml', '.yml', '.csv']
        self.max_file_size = 10 * 1024 * 1024  # 10MB limit
        self.allowed_directories = set()
        
        # Ensure base path exists
        self.base_path.mkdir(parents=True, exist_ok=True)
        
        logger.info(f"File system tool initialized with base path: {self.base_path}")
    
    def _validate_path(self, file_path: str) -> Path:
        """
        Validate and sanitize file path for security
        
        Args:
            file_path: Relative file path
            
        Returns:
            Resolved Path object
            
        Raises:
            ValueError: If path is invalid or outside allowed directory
        """
        # Convert to Path object
        path = Path(file_path)
        
        # Resolve relative to base path
        resolved_path = (self.base_path / path).resolve()
        
        # Check if path is within base directory (security check)
        try:
            resolved_path.relative_to(self.base_path)
        except ValueError:
            raise ValueError(f"Path {file_path} is outside allowed directory")
        
        return resolved_path
    
    def _check_file_size(self, file_path: Path) -> bool:
        """Check if file size is within limits"""
        try:
            return file_path.stat().st_size <= self.max_file_size
        except OSError:
            return False
    
    def _get_file_hash(self, file_path: Path) -> str:
        """Generate SHA-256 hash of file content"""
        try:
            with open(file_path, 'rb') as f:
                content = f.read()
                return hashlib.sha256(content).hexdigest()
        except Exception as e:
            logger.error(f"Failed to generate file hash: {e}")
            return ""
    
    def _is_allowed_extension(self, file_path: Path) -> bool:
        """Check if file extension is allowed"""
        return file_path.suffix.lower() in self.allowed_extensions
    
    def read_file(self, file_path: str) -> FileOperation:
        """
        Read file content with security validation
        
        Args:
            file_path: Relative path to file
            
        Returns:
            FileOperation with file content
        """
        try:
            path = self._validate_path(file_path)
            
            # Check if file exists
            if not path.exists():
                return FileOperation(
                    success=False,
                    operation="read",
                    file_path=str(path),
                    error="File does not exist"
                )
            
            # Check if it's a file (not directory)
            if not path.is_file():
                return FileOperation(
                    success=False,
                    operation="read",
                    file_path=str(path),
                    error="Path is not a file"
                )
            
            # Check file size
            if not self._check_file_size(path):
                return FileOperation(
                    success=False,
                    operation="read",
                    file_path=str(path),
                    error="File size exceeds limit"
                )
            
            # Check file extension
            if not self._is_allowed_extension(path):
                return FileOperation(
                    success=False,
                    operation="read",
                    file_path=str(path),
                    error="File type not allowed"
                )
            
            # Read file content
            with open(path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            file_hash = self._get_file_hash(path)
            
            logger.info(f"Successfully read file: {path}")
            
            return FileOperation(
                success=True,
                operation="read",
                file_path=str(path),
                data=content,
                file_hash=file_hash
            )
            
        except Exception as e:
            logger.error(f"Failed to read file {file_path}: {e}")
            return FileOperation(
                success=False,
                operation="read",
                file_path=file_path,
                error=str(e)
            )
    
    def write_file(self, file_path: str, content: str, overwrite: bool = False) -> FileOperation:
        """
        Write content to file with security validation
        
        Args:
            file_path: Relative path to file
            content: Content to write
            overwrite: Whether to overwrite existing file
            
        Returns:
            FileOperation with result
        """
        try:
            path = self._validate_path(file_path)
            
            # Check file extension
            if not self._is_allowed_extension(path):
                return FileOperation(
                    success=False,
                    operation="write",
                    file_path=str(path),
                    error="File type not allowed"
                )
            
            # Check if file exists and overwrite is not allowed
            if path.exists() and not overwrite:
                return FileOperation(
                    success=False,
                    operation="write",
                    file_path=str(path),
                    error="File already exists and overwrite not allowed"
                )
            
            # Ensure directory exists
            path.parent.mkdir(parents=True, exist_ok=True)
            
            # Write content
            with open(path, 'w', encoding='utf-8') as f:
                f.write(content)
            
            file_hash = self._get_file_hash(path)
            
            logger.info(f"Successfully wrote file: {path}")
            
            return FileOperation(
                success=True,
                operation="write",
                file_path=str(path),
                data=content,
                file_hash=file_hash
            )
            
        except Exception as e:
            logger.error(f"Failed to write file {file_path}: {e}")
            return FileOperation(
                success=False,
                operation="write",
                file_path=file_path,
                error=str(e)
            )
    
    def read_json(self, file_path: str) -> FileOperation:
        """
        Read and parse JSON file
        
        Args:
            file_path: Relative path to JSON file
            
        Returns:
            FileOperation with parsed JSON data
        """
        read_result = self.read_file(file_path)
        
        if not read_result.success:
            return read_result
        
        try:
            json_data = json.loads(read_result.data)
            return FileOperation(
                success=True,
                operation="read_json",
                file_path=read_result.file_path,
                data=json_data,
                file_hash=read_result.file_hash
            )
        except json.JSONDecodeError as e:
            return FileOperation(
                success=False,
                operation="read_json",
                file_path=file_path,
                error=f"Invalid JSON: {e}"
            )
    
    def write_json(self, file_path: str, data: Dict[str, Any], 
                  indent: int = 2, overwrite: bool = False) -> FileOperation:
        """
        Write data as JSON file
        
        Args:
            file_path: Relative path to JSON file
            data: Data to write as JSON
            indent: JSON indentation
            overwrite: Whether to overwrite existing file
            
        Returns:
            FileOperation with result
        """
        try:
            json_content = json.dumps(data, indent=indent, ensure_ascii=False)
            return self.write_file(file_path, json_content, overwrite)
        except Exception as e:
            return FileOperation(
                success=False,
                operation="write_json",
                file_path=file_path,
                error=f"Failed to serialize JSON: {e}"
            )
    
    def list_directory(self, directory_path: str = ".") -> FileOperation:
        """
        List files and directories
        
        Args:
            directory_path: Relative path to directory
            
        Returns:
            FileOperation with directory listing
        """
        try:
            path = self._validate_path(directory_path)
            
            if not path.exists():
                return FileOperation(
                    success=False,
                    operation="list",
                    file_path=str(path),
                    error="Directory does not exist"
                )
            
            if not path.is_dir():
                return FileOperation(
                    success=False,
                    operation="list",
                    file_path=str(path),
                    error="Path is not a directory"
                )
            
            files = []
            for item in path.iterdir():
                try:
                    stat = item.stat()
                    files.append(FileInfo(
                        path=str(item.relative_to(self.base_path)),
                        name=item.name,
                        size=stat.st_size,
                        modified=datetime.fromtimestamp(stat.st_mtime),
                        mime_type=mimetypes.guess_type(item.name)[0] or "unknown",
                        is_directory=item.is_dir(),
                        permissions=oct(stat.st_mode)[-3:]
                    ))
                except OSError:
                    # Skip files we can't access
                    continue
            
            # Sort by name
            files.sort(key=lambda x: x.name)
            
            logger.info(f"Successfully listed directory: {path}")
            
            return FileOperation(
                success=True,
                operation="list",
                file_path=str(path),
                data=files
            )
            
        except Exception as e:
            logger.error(f"Failed to list directory {directory_path}: {e}")
            return FileOperation(
                success=False,
                operation="list",
                file_path=directory_path,
                error=str(e)
            )
    
    def get_file_info(self, file_path: str) -> FileOperation:
        """
        Get detailed file information
        
        Args:
            file_path: Relative path to file
            
        Returns:
            FileOperation with file info
        """
        try:
            path = self._validate_path(file_path)
            
            if not path.exists():
                return FileOperation(
                    success=False,
                    operation="info",
                    file_path=str(path),
                    error="File does not exist"
                )
            
            stat = path.stat()
            file_info = FileInfo(
                path=str(path.relative_to(self.base_path)),
                name=path.name,
                size=stat.st_size,
                modified=datetime.fromtimestamp(stat.st_mtime),
                mime_type=mimetypes.guess_type(path.name)[0] or "unknown",
                is_directory=path.is_dir(),
                permissions=oct(stat.st_mode)[-3:]
            )
            
            logger.info(f"Successfully got file info: {path}")
            
            return FileOperation(
                success=True,
                operation="info",
                file_path=str(path),
                data=file_info
            )
            
        except Exception as e:
            logger.error(f"Failed to get file info {file_path}: {e}")
            return FileOperation(
                success=False,
                operation="info",
                file_path=file_path,
                error=str(e)
            )
    
    def copy_file(self, source_path: str, dest_path: str, overwrite: bool = False) -> FileOperation:
        """
        Copy file from source to destination
        
        Args:
            source_path: Source file path
            dest_path: Destination file path
            overwrite: Whether to overwrite existing file
            
        Returns:
            FileOperation with result
        """
        try:
            source = self._validate_path(source_path)
            dest = self._validate_path(dest_path)
            
            if not source.exists():
                return FileOperation(
                    success=False,
                    operation="copy",
                    file_path=str(source),
                    error="Source file does not exist"
                )
            
            if not source.is_file():
                return FileOperation(
                    success=False,
                    operation="copy",
                    file_path=str(source),
                    error="Source is not a file"
                )
            
            if dest.exists() and not overwrite:
                return FileOperation(
                    success=False,
                    operation="copy",
                    file_path=str(dest),
                    error="Destination file exists and overwrite not allowed"
                )
            
            # Ensure destination directory exists
            dest.parent.mkdir(parents=True, exist_ok=True)
            
            # Copy file
            shutil.copy2(source, dest)
            
            logger.info(f"Successfully copied file: {source} -> {dest}")
            
            return FileOperation(
                success=True,
                operation="copy",
                file_path=str(dest)
            )
            
        except Exception as e:
            logger.error(f"Failed to copy file {source_path} -> {dest_path}: {e}")
            return FileOperation(
                success=False,
                operation="copy",
                file_path=f"{source_path} -> {dest_path}",
                error=str(e)
            )

def main():
    """Test the file system tool"""
    fs_tool = FileSystemTool("test_files")
    
    # Test writing a file
    result = fs_tool.write_file("test.txt", "Hello, World!")
    print(f"Write result: {result.success}")
    
    # Test reading the file
    result = fs_tool.read_file("test.txt")
    print(f"Read result: {result.success}")
    if result.success:
        print(f"Content: {result.data}")
    
    # Test writing JSON
    data = {"name": "test", "value": 123}
    result = fs_tool.write_json("test.json", data)
    print(f"JSON write result: {result.success}")
    
    # Test reading JSON
    result = fs_tool.read_json("test.json")
    print(f"JSON read result: {result.success}")
    if result.success:
        print(f"JSON data: {result.data}")
    
    # Test listing directory
    result = fs_tool.list_directory(".")
    print(f"List result: {result.success}")
    if result.success:
        for file_info in result.data:
            print(f"- {file_info.name} ({file_info.size} bytes)")

if __name__ == "__main__":
    main()
