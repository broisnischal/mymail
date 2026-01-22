export const config = {
  // Server
  api: {
    port: parseInt(process.env.API_PORT || '3000'),
    host: process.env.API_HOST || '0.0.0.0',
    jwtSecret: process.env.JWT_SECRET || 'change-me-in-production',
    jwtExpiry: process.env.JWT_EXPIRY || '7d',
  },
   
  // Database
  database: {
    url: process.env.DATABASE_URL || 'postgresql://postgres:postgres@postgres:5432/mymail',
    maxConnections: parseInt(process.env.DB_MAX_CONNECTIONS || '10'),
  },
  
  // Redis
  redis: {
    url: process.env.REDIS_URL || 'redis://redis:6379',
    prefix: process.env.REDIS_PREFIX || 'mymail:',
  },
  
  // MinIO
  minio: {
    endpoint: process.env.MINIO_ENDPOINT || 'minio:9000',
    accessKey: process.env.MINIO_ACCESS_KEY || 'minioadmin',
    secretKey: process.env.MINIO_SECRET_KEY || 'minioadmin',
    bucket: process.env.MINIO_BUCKET || 'mails',
    useSSL: process.env.MINIO_USE_SSL === 'true',
  },
  
  // SMTP
  smtp: {
    port: parseInt(process.env.SMTP_PORT || '25'),
    host: process.env.SMTP_HOST || '0.0.0.0',
    domain: process.env.SMTP_DOMAIN || 'mymail.com',
    maxMessageSize: parseInt(process.env.SMTP_MAX_SIZE || '10485760'), // 10MB
  },
  
  // Worker
  worker: {
    concurrency: parseInt(process.env.WORKER_CONCURRENCY || '10'),
    batchSize: parseInt(process.env.WORKER_BATCH_SIZE || '100'),
  },
  
  // Rate Limiting
  rateLimit: {
    emailsPerUser: parseInt(process.env.RATE_LIMIT_EMAILS_PER_USER || '1000'),
    emailsPerHour: parseInt(process.env.RATE_LIMIT_EMAILS_PER_HOUR || '100'),
    connectionsPerIP: parseInt(process.env.RATE_LIMIT_CONNECTIONS_PER_IP || '10'),
  },
  
  // DKIM
  dkim: {
    privateKey: process.env.DKIM_PRIVATE_KEY || '',
    selector: process.env.DKIM_SELECTOR || 'default',
    domain: process.env.DKIM_DOMAIN || process.env.SMTP_DOMAIN || 'mymail.com',
  },
  
  // Temp Mail
  tempMail: {
    enabled: process.env.TEMP_MAIL_ENABLED !== 'false',
    ttl: parseInt(process.env.TEMP_MAIL_TTL || '86400'), // 24 hours
  },
};
