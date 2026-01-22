import { Client } from 'minio';
import { config } from '@shared/config';

export const minioClient = new Client({
  endPoint: config.minio.endpoint.split(':')[0],
  port: parseInt(config.minio.endpoint.split(':')[1] || '9000'),
  useSSL: config.minio.useSSL,
  accessKey: config.minio.accessKey,
  secretKey: config.minio.secretKey,
});

export async function ensureBucket() {
  const exists = await minioClient.bucketExists(config.minio.bucket);
  if (!exists) {
    await minioClient.makeBucket(config.minio.bucket, 'us-east-1');
  }
}

export async function uploadEmail(path: string, content: Buffer): Promise<string> {
  await ensureBucket();
  await minioClient.putObject(config.minio.bucket, path, content);
  return path;
}

export async function getEmail(path: string): Promise<Buffer> {
  const stream = await minioClient.getObject(config.minio.bucket, path);
  const chunks: Buffer[] = [];
  
  return new Promise((resolve, reject) => {
    stream.on('data', (chunk) => chunks.push(chunk));
    stream.on('end', () => resolve(Buffer.concat(chunks)));
    stream.on('error', reject);
  });
}

export async function deleteEmail(path: string): Promise<void> {
  await minioClient.removeObject(config.minio.bucket, path);
}
