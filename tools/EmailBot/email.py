#!/usr/bin/env python3
"""
Email Tool for CMP Agents
Provides email sending and reading capabilities with security controls
"""

import smtplib
import imaplib
import email
import logging
from typing import Dict, Any, Optional, List
from dataclasses import dataclass
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from email.mime.base import MIMEBase
from email import encoders
from datetime import datetime
import ssl
import re

logger = logging.getLogger(__name__)

@dataclass
class EmailMessage:
    """Represents an email message"""
    subject: str
    sender: str
    recipients: List[str]
    body: str
    html_body: Optional[str] = None
    attachments: List[str] = None
    message_id: Optional[str] = None
    date: Optional[datetime] = None

@dataclass
class EmailConfig:
    """Email configuration"""
    smtp_server: str
    smtp_port: int
    imap_server: str
    imap_port: int
    username: str
    password: str
    use_ssl: bool = True
    use_tls: bool = True

@dataclass
class EmailResult:
    """Represents email operation result"""
    success: bool
    operation: str
    message: str
    data: Optional[Any] = None
    error: Optional[str] = None

class EmailTool:
    """Email tool for agents with security and spam protection"""
    
    def __init__(self, config: EmailConfig):
        self.config = config
        self.max_recipients = 10
        self.max_attachment_size = 10 * 1024 * 1024  # 10MB
        self.allowed_domains = set()
        self.blocked_domains = set()
        
        # Security patterns
        self.suspicious_patterns = [
            r'password.*reset',
            r'click.*here',
            r'urgent.*action',
            r'account.*suspended',
            r'verify.*email'
        ]
        
        logger.info(f"Email tool initialized for {config.username}")
    
    def _validate_email_address(self, email_addr: str) -> bool:
        """Validate email address format"""
        pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
        return bool(re.match(pattern, email_addr))
    
    def _validate_domain(self, email_addr: str) -> bool:
        """Check if email domain is allowed/blocked"""
        domain = email_addr.split('@')[1].lower()
        
        if domain in self.blocked_domains:
            return False
        
        if self.allowed_domains and domain not in self.allowed_domains:
            return False
        
        return True
    
    def _check_suspicious_content(self, subject: str, body: str) -> bool:
        """Check for suspicious content patterns"""
        content = f"{subject} {body}".lower()
        
        for pattern in self.suspicious_patterns:
            if re.search(pattern, content):
                logger.warning(f"Suspicious content detected: {pattern}")
                return True
        
        return False
    
    def _sanitize_content(self, content: str) -> str:
        """Sanitize email content for security"""
        # Remove potentially dangerous HTML
        content = re.sub(r'<script.*?</script>', '', content, flags=re.IGNORECASE | re.DOTALL)
        content = re.sub(r'<iframe.*?</iframe>', '', content, flags=re.IGNORECASE | re.DOTALL)
        content = re.sub(r'javascript:', '', content, flags=re.IGNORECASE)
        
        return content
    
    def send_email(self, message: EmailMessage) -> EmailResult:
        """
        Send email with security validation
        
        Args:
            message: EmailMessage object with email details
            
        Returns:
            EmailResult with operation result
        """
        try:
            # Validate recipients
            if len(message.recipients) > self.max_recipients:
                return EmailResult(
                    success=False,
                    operation="send",
                    message="Too many recipients",
                    error=f"Maximum {self.max_recipients} recipients allowed"
                )
            
            for recipient in message.recipients:
                if not self._validate_email_address(recipient):
                    return EmailResult(
                        success=False,
                        operation="send",
                        message="Invalid email address",
                        error=f"Invalid email format: {recipient}"
                    )
                
                if not self._validate_domain(recipient):
                    return EmailResult(
                        success=False,
                        operation="send",
                        message="Domain not allowed",
                        error=f"Domain not allowed for: {recipient}"
                    )
            
            # Check for suspicious content
            if self._check_suspicious_content(message.subject, message.body):
                return EmailResult(
                    success=False,
                    operation="send",
                    message="Suspicious content detected",
                    error="Email contains suspicious patterns"
                )
            
            # Sanitize content
            sanitized_subject = self._sanitize_content(message.subject)
            sanitized_body = self._sanitize_content(message.body)
            
            # Create email message
            msg = MIMEMultipart()
            msg['From'] = self.config.username
            msg['To'] = ', '.join(message.recipients)
            msg['Subject'] = sanitized_subject
            
            # Add text body
            msg.attach(MIMEText(sanitized_body, 'plain'))
            
            # Add HTML body if provided
            if message.html_body:
                sanitized_html = self._sanitize_content(message.html_body)
                msg.attach(MIMEText(sanitized_html, 'html'))
            
            # Add attachments
            if message.attachments:
                for attachment_path in message.attachments:
                    try:
                        with open(attachment_path, 'rb') as attachment:
                            part = MIMEBase('application', 'octet-stream')
                            part.set_payload(attachment.read())
                        
                        encoders.encode_base64(part)
                        part.add_header(
                            'Content-Disposition',
                            f'attachment; filename= {attachment_path.split("/")[-1]}'
                        )
                        msg.attach(part)
                    except Exception as e:
                        logger.error(f"Failed to attach file {attachment_path}: {e}")
            
            # Send email
            if self.config.use_ssl:
                server = smtplib.SMTP_SSL(self.config.smtp_server, self.config.smtp_port)
            else:
                server = smtplib.SMTP(self.config.smtp_server, self.config.smtp_port)
                if self.config.use_tls:
                    server.starttls()
            
            server.login(self.config.username, self.config.password)
            server.send_message(msg)
            server.quit()
            
            logger.info(f"Email sent successfully to {message.recipients}")
            
            return EmailResult(
                success=True,
                operation="send",
                message="Email sent successfully",
                data={"recipients": message.recipients, "subject": sanitized_subject}
            )
            
        except Exception as e:
            logger.error(f"Failed to send email: {e}")
            return EmailResult(
                success=False,
                operation="send",
                message="Failed to send email",
                error=str(e)
            )
    
    def read_emails(self, folder: str = "INBOX", limit: int = 10, 
                   unread_only: bool = False) -> EmailResult:
        """
        Read emails from specified folder
        
        Args:
            folder: Email folder to read from
            limit: Maximum number of emails to read
            unread_only: Only read unread emails
            
        Returns:
            EmailResult with email list
        """
        try:
            # Connect to IMAP server
            if self.config.use_ssl:
                mail = imaplib.IMAP4_SSL(self.config.imap_server, self.config.imap_port)
            else:
                mail = imaplib.IMAP4(self.config.imap_server, self.config.imap_port)
            
            mail.login(self.config.username, self.config.password)
            mail.select(folder)
            
            # Search criteria
            if unread_only:
                status, messages = mail.search(None, 'UNSEEN')
            else:
                status, messages = mail.search(None, 'ALL')
            
            if status != 'OK':
                return EmailResult(
                    success=False,
                    operation="read",
                    message="Failed to search emails",
                    error="IMAP search failed"
                )
            
            # Get email IDs
            email_ids = messages[0].split()
            
            # Limit number of emails
            if len(email_ids) > limit:
                email_ids = email_ids[-limit:]
            
            emails = []
            
            for email_id in email_ids:
                try:
                    status, msg_data = mail.fetch(email_id, '(RFC822)')
                    
                    if status != 'OK':
                        continue
                    
                    raw_email = msg_data[0][1]
                    email_message = email.message_from_bytes(raw_email)
                    
                    # Extract email details
                    subject = email.header.decode_header(email_message['subject'])[0][0]
                    if isinstance(subject, bytes):
                        subject = subject.decode('utf-8', errors='ignore')
                    
                    sender = email_message['from']
                    date_str = email_message['date']
                    date = email.utils.parsedate_to_datetime(date_str) if date_str else None
                    
                    # Extract body
                    body = ""
                    html_body = ""
                    
                    if email_message.is_multipart():
                        for part in email_message.walk():
                            content_type = part.get_content_type()
                            content_disposition = str(part.get('Content-Disposition'))
                            
                            if "attachment" not in content_disposition:
                                if content_type == "text/plain":
                                    body = part.get_payload(decode=True).decode('utf-8', errors='ignore')
                                elif content_type == "text/html":
                                    html_body = part.get_payload(decode=True).decode('utf-8', errors='ignore')
                    else:
                        body = email_message.get_payload(decode=True).decode('utf-8', errors='ignore')
                    
                    emails.append({
                        'id': email_id.decode(),
                        'subject': subject,
                        'sender': sender,
                        'date': date,
                        'body': body,
                        'html_body': html_body,
                        'unread': '\\Seen' not in mail.fetch(email_id, '(FLAGS)')[1][0].decode()
                    })
                    
                except Exception as e:
                    logger.error(f"Failed to parse email {email_id}: {e}")
                    continue
            
            mail.logout()
            
            logger.info(f"Successfully read {len(emails)} emails from {folder}")
            
            return EmailResult(
                success=True,
                operation="read",
                message=f"Read {len(emails)} emails",
                data=emails
            )
            
        except Exception as e:
            logger.error(f"Failed to read emails: {e}")
            return EmailResult(
                success=False,
                operation="read",
                message="Failed to read emails",
                error=str(e)
            )
    
    def mark_as_read(self, email_ids: List[str], folder: str = "INBOX") -> EmailResult:
        """
        Mark emails as read
        
        Args:
            email_ids: List of email IDs to mark as read
            folder: Email folder
            
        Returns:
            EmailResult with operation result
        """
        try:
            mail = imaplib.IMAP4_SSL(self.config.imap_server, self.config.imap_port)
            mail.login(self.config.username, self.config.password)
            mail.select(folder)
            
            for email_id in email_ids:
                mail.store(email_id, '+FLAGS', '\\Seen')
            
            mail.logout()
            
            logger.info(f"Marked {len(email_ids)} emails as read")
            
            return EmailResult(
                success=True,
                operation="mark_read",
                message=f"Marked {len(email_ids)} emails as read"
            )
            
        except Exception as e:
            logger.error(f"Failed to mark emails as read: {e}")
            return EmailResult(
                success=False,
                operation="mark_read",
                message="Failed to mark emails as read",
                error=str(e)
            )
    
    def delete_emails(self, email_ids: List[str], folder: str = "INBOX") -> EmailResult:
        """
        Delete emails
        
        Args:
            email_ids: List of email IDs to delete
            folder: Email folder
            
        Returns:
            EmailResult with operation result
        """
        try:
            mail = imaplib.IMAP4_SSL(self.config.imap_server, self.config.imap_port)
            mail.login(self.config.username, self.config.password)
            mail.select(folder)
            
            for email_id in email_ids:
                mail.store(email_id, '+FLAGS', '\\Deleted')
            
            mail.expunge()
            mail.logout()
            
            logger.info(f"Deleted {len(email_ids)} emails")
            
            return EmailResult(
                success=True,
                operation="delete",
                message=f"Deleted {len(email_ids)} emails"
            )
            
        except Exception as e:
            logger.error(f"Failed to delete emails: {e}")
            return EmailResult(
                success=False,
                operation="delete",
                message="Failed to delete emails",
                error=str(e)
            )
    
    def get_folders(self) -> EmailResult:
        """
        Get list of email folders
        
        Returns:
            EmailResult with folder list
        """
        try:
            mail = imaplib.IMAP4_SSL(self.config.imap_server, self.config.imap_port)
            mail.login(self.config.username, self.config.password)
            
            status, folders = mail.list()
            
            if status != 'OK':
                return EmailResult(
                    success=False,
                    operation="list_folders",
                    message="Failed to list folders",
                    error="IMAP list failed"
                )
            
            folder_list = []
            for folder in folders:
                folder_name = folder.decode().split('"')[-2]
                folder_list.append(folder_name)
            
            mail.logout()
            
            logger.info(f"Found {len(folder_list)} email folders")
            
            return EmailResult(
                success=True,
                operation="list_folders",
                message=f"Found {len(folder_list)} folders",
                data=folder_list
            )
            
        except Exception as e:
            logger.error(f"Failed to list folders: {e}")
            return EmailResult(
                success=False,
                operation="list_folders",
                message="Failed to list folders",
                error=str(e)
            )

def main():
    """Test the email tool"""
    # Example configuration (replace with actual values)
    config = EmailConfig(
        smtp_server="smtp.gmail.com",
        smtp_port=587,
        imap_server="imap.gmail.com",
        imap_port=993,
        username="your_email@gmail.com",
        password="your_app_password",
        use_ssl=True,
        use_tls=True
    )
    
    email_tool = EmailTool(config)
    
    # Test sending email
    message = EmailMessage(
        subject="Test Email from CMP Agent",
        sender=config.username,
        recipients=["recipient@example.com"],
        body="This is a test email sent by the CMP agent.",
        html_body="<h1>Test Email</h1><p>This is a test email sent by the CMP agent.</p>"
    )
    
    # Uncomment to test sending
    # result = email_tool.send_email(message)
    # print(f"Send result: {result.success}")
    # if not result.success:
    #     print(f"Error: {result.error}")
    
    # Test reading emails
    # result = email_tool.read_emails(limit=5, unread_only=True)
    # print(f"Read result: {result.success}")
    # if result.success:
    #     for email_data in result.data:
    #         print(f"- {email_data['subject']} from {email_data['sender']}")
    
    # Test getting folders
    # result = email_tool.get_folders()
    # print(f"Folders result: {result.success}")
    # if result.success:
    #     print(f"Folders: {result.data}")

if __name__ == "__main__":
    main()
