#!/usr/bin/env python3
"""
Send test email using Python's built-in smtplib
No additional dependencies required!
"""

import smtplib
from email.mime.text import MIMEText
from datetime import datetime
import sys

def send_test_email(to_email, from_email="sender@example.com", smtp_host="localhost", smtp_port=2525):
    """Send a test email via SMTP"""
    
    # Create message
    msg = MIMEText(f"""This is a test email sent via Python smtplib to test the MyMail SMTP server.

Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}
This email was sent using the send-test-email-python.py script.
""")
    
    msg['Subject'] = f"Test Email - {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}"
    msg['From'] = from_email
    msg['To'] = to_email
    
    try:
        # Connect to SMTP server
        print(f"Connecting to {smtp_host}:{smtp_port}...")
        server = smtplib.SMTP(smtp_host, smtp_port)
        server.set_debuglevel(1)  # Show SMTP conversation
        
        # Send email
        print(f"Sending email from {from_email} to {to_email}...")
        server.send_message(msg)
        server.quit()
        
        print("✓ Email sent successfully!")
        return True
        
    except Exception as e:
        print(f"✗ Error sending email: {e}")
        return False

if __name__ == "__main__":
    to_email = sys.argv[1] if len(sys.argv) > 1 else "test@mail.localhost"
    from_email = sys.argv[2] if len(sys.argv) > 2 else "sender@example.com"
    
    print(f"=== Sending Test Email ===")
    print(f"To: {to_email}")
    print(f"From: {from_email}")
    print()
    
    send_test_email(to_email, from_email)
