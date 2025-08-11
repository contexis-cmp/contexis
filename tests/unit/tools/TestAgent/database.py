#!/usr/bin/env python3
"""
Database Tool for CMP Agents
Provides database query capabilities with connection pooling and security
"""

import sqlite3
import json
import logging
from typing import List, Dict, Optional, Any
from dataclasses import dataclass
from contextlib import contextmanager
import hashlib
import os

logger = logging.getLogger(__name__)

@dataclass
class QueryResult:
    """Represents a database query result"""
    data: List[Dict[str, Any]]
    row_count: int
    columns: List[str]
    execution_time: float
    query_hash: str

class DatabaseTool:
    """Database tool for agents with security and connection management"""
    
    def __init__(self, db_path: str, max_connections: int = 5):
        self.db_path = db_path
        self.max_connections = max_connections
        self.allowed_tables = set()
        self.query_whitelist = set()
        
        # Initialize database if it doesn't exist
        self._init_database()
        
    def _init_database(self):
        """Initialize database with required tables"""
        try:
            with self._get_connection() as conn:
                cursor = conn.cursor()
                
                # Create users table if it doesn't exist
                cursor.execute("""
                    CREATE TABLE IF NOT EXISTS users (
                        id INTEGER PRIMARY KEY,
                        username TEXT UNIQUE NOT NULL,
                        email TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        status TEXT DEFAULT 'active'
                    )
                """)
                
                # Create orders table if it doesn't exist
                cursor.execute("""
                    CREATE TABLE IF NOT EXISTS orders (
                        id INTEGER PRIMARY KEY,
                        user_id INTEGER,
                        order_number TEXT UNIQUE NOT NULL,
                        total_amount REAL,
                        status TEXT DEFAULT 'pending',
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        FOREIGN KEY (user_id) REFERENCES users (id)
                    )
                """)
                
                # Create products table if it doesn't exist
                cursor.execute("""
                    CREATE TABLE IF NOT EXISTS products (
                        id INTEGER PRIMARY KEY,
                        name TEXT NOT NULL,
                        description TEXT,
                        price REAL,
                        stock_quantity INTEGER DEFAULT 0,
                        category TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                    )
                """)
                
                conn.commit()
                logger.info("Database initialized successfully")
                
        except Exception as e:
            logger.error(f"Failed to initialize database: {e}")
            raise
    
    @contextmanager
    def _get_connection(self):
        """Get database connection with proper error handling"""
        conn = None
        try:
            conn = sqlite3.connect(self.db_path, timeout=30.0)
            conn.row_factory = sqlite3.Row  # Enable dict-like access
            yield conn
        except Exception as e:
            logger.error(f"Database connection error: {e}")
            raise
        finally:
            if conn:
                conn.close()
    
    def _validate_query(self, query: str) -> bool:
        """
        Validate query for security and allowed operations
        
        Args:
            query: SQL query to validate
            
        Returns:
            True if query is allowed, False otherwise
        """
        query_lower = query.lower().strip()
        
        # Block dangerous operations
        dangerous_keywords = [
            'drop', 'delete', 'truncate', 'alter', 'create', 'insert', 'update'
        ]
        
        for keyword in dangerous_keywords:
            if keyword in query_lower:
                logger.warning(f"Blocked query with dangerous keyword: {keyword}")
                return False
        
        # Only allow SELECT queries for safety
        if not query_lower.startswith('select'):
            logger.warning("Only SELECT queries are allowed")
            return False
        
        return True
    
    def _hash_query(self, query: str) -> str:
        """Generate hash for query tracking"""
        return hashlib.sha256(query.encode()).hexdigest()[:16]
    
    def query(self, sql: str, params: Optional[Dict[str, Any]] = None) -> QueryResult:
        """
        Execute a safe database query
        
        Args:
            sql: SQL query string (SELECT only)
            params: Query parameters for prepared statements
            
        Returns:
            QueryResult object with data and metadata
        """
        import time
        
        start_time = time.time()
        
        try:
            # Validate query
            if not self._validate_query(sql):
                raise ValueError("Query validation failed")
            
            # Generate query hash
            query_hash = self._hash_query(sql)
            
            logger.info(f"Executing query: {sql[:100]}...")
            
            with self._get_connection() as conn:
                cursor = conn.cursor()
                
                # Execute query with parameters if provided
                if params:
                    cursor.execute(sql, params)
                else:
                    cursor.execute(sql)
                
                # Fetch results
                rows = cursor.fetchall()
                columns = [description[0] for description in cursor.description] if cursor.description else []
                
                # Convert to list of dictionaries
                data = []
                for row in rows:
                    data.append(dict(row))
                
                execution_time = time.time() - start_time
                
                result = QueryResult(
                    data=data,
                    row_count=len(data),
                    columns=columns,
                    execution_time=execution_time,
                    query_hash=query_hash
                )
                
                logger.info(f"Query completed: {result.row_count} rows in {execution_time:.3f}s")
                return result
                
        except Exception as e:
            logger.error(f"Query execution failed: {e}")
            raise
    
    def get_user_info(self, user_id: int) -> Optional[Dict[str, Any]]:
        """
        Get user information by ID
        
        Args:
            user_id: User ID to lookup
            
        Returns:
            User information dictionary or None
        """
        try:
            result = self.query(
                "SELECT id, username, email, status, created_at FROM users WHERE id = ?",
                {"user_id": user_id}
            )
            
            if result.data:
                return result.data[0]
            return None
            
        except Exception as e:
            logger.error(f"Failed to get user info: {e}")
            return None
    
    def get_user_orders(self, user_id: int) -> List[Dict[str, Any]]:
        """
        Get orders for a specific user
        
        Args:
            user_id: User ID to get orders for
            
        Returns:
            List of order dictionaries
        """
        try:
            result = self.query(
                """
                SELECT o.id, o.order_number, o.total_amount, o.status, o.created_at
                FROM orders o
                WHERE o.user_id = ?
                ORDER BY o.created_at DESC
                """,
                {"user_id": user_id}
            )
            
            return result.data
            
        except Exception as e:
            logger.error(f"Failed to get user orders: {e}")
            return []
    
    def search_products(self, search_term: str, category: Optional[str] = None) -> List[Dict[str, Any]]:
        """
        Search products by name or description
        
        Args:
            search_term: Search term to match against product name/description
            category: Optional category filter
            
        Returns:
            List of matching products
        """
        try:
            if category:
                result = self.query(
                    """
                    SELECT id, name, description, price, stock_quantity, category
                    FROM products
                    WHERE (name LIKE ? OR description LIKE ?) AND category = ?
                    ORDER BY name
                    """,
                    {
                        "search_term": f"%{search_term}%",
                        "category": category
                    }
                )
            else:
                result = self.query(
                    """
                    SELECT id, name, description, price, stock_quantity, category
                    FROM products
                    WHERE name LIKE ? OR description LIKE ?
                    ORDER BY name
                    """,
                    {"search_term": f"%{search_term}%"}
                )
            
            return result.data
            
        except Exception as e:
            logger.error(f"Failed to search products: {e}")
            return []
    
    def get_database_stats(self) -> Dict[str, Any]:
        """
        Get database statistics
        
        Returns:
            Dictionary with database statistics
        """
        try:
            stats = {}
            
            # Get table counts
            tables = ['users', 'orders', 'products']
            for table in tables:
                result = self.query(f"SELECT COUNT(*) as count FROM {table}")
                if result.data:
                    stats[f"{table}_count"] = result.data[0]['count']
            
            # Get database size
            if os.path.exists(self.db_path):
                stats['database_size_mb'] = os.path.getsize(self.db_path) / (1024 * 1024)
            
            return stats
            
        except Exception as e:
            logger.error(f"Failed to get database stats: {e}")
            return {}

def main():
    """Test the database tool"""
    tool = DatabaseTool("test_agent.db")
    
    # Test basic query
    result = tool.query("SELECT name FROM sqlite_master WHERE type='table'")
    print("Database tables:")
    for row in result.data:
        print(f"- {row['name']}")
    
    # Test user lookup
    user_info = tool.get_user_info(1)
    if user_info:
        print(f"\nUser info: {user_info}")
    
    # Test product search
    products = tool.search_products("test")
    print(f"\nFound {len(products)} products")
    
    # Test database stats
    stats = tool.get_database_stats()
    print(f"\nDatabase stats: {stats}")

if __name__ == "__main__":
    main()
