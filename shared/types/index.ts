export interface User {
  id: string;
  email: string;
  passwordHash: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface Mailbox {
  id: string;
  userId: string;
  address: string;
  isAlias: boolean;
  isTemp: boolean;
  createdAt: Date;
  updatedAt: Date;
}

export interface Email {
  id: string;
  mailboxId: string;
  messageId: string;
  from: string;
  to: string[];
  cc?: string[];
  bcc?: string[];
  subject: string;
  textBody?: string;
  htmlBody?: string;
  minioPath: string;
  size: number;
  receivedAt: Date;
  createdAt: Date;
}

export interface EmailMetadata {
  id: string;
  emailId: string;
  headers: Record<string, string>;
  attachments?: Array<{
    filename: string;
    contentType: string;
    size: number;
    minioPath: string;
  }>;
}

export interface QueueJob {
  id: string;
  type: 'process_email' | 'send_email' | 'cleanup_temp';
  payload: Record<string, any>;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  attempts: number;
  createdAt: Date;
  processedAt?: Date;
}

export interface AuthToken {
  userId: string;
  token: string;
  expiresAt: Date;
}
