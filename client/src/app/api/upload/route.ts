import { NextResponse, NextRequest } from 'next/server';

export interface UploadResponse {
  torrentID: string;
}

export async function POST(req: NextRequest): Promise<NextResponse> {
  const formData = await req.formData();
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/upload`, {
    method: 'POST',
    body: formData,
  });
  const data: UploadResponse = await response.json();
  return NextResponse.json(data);
}
