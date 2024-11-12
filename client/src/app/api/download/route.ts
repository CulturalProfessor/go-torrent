import { NextResponse } from 'next/server';

export async function GET(req: Request) {
  const { searchParams } = new URL(req.url);
  const torrentID = searchParams.get('id');

  if (!torrentID) {
    return NextResponse.json({ error: 'Torrent ID is required' }, { status: 400 });
  }

  try {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}//download?id=${torrentID}`);

    if (!response.ok) {
      return NextResponse.json({ error: 'File not found or download incomplete' }, { status: 404 });
    }

    return new NextResponse(response.body, {
      headers: {
        'Content-Disposition': response.headers.get('Content-Disposition') || 'attachment; filename="downloaded-file"',
        'Content-Type': response.headers.get('Content-Type') || 'application/octet-stream',
      },
    });
  } catch (error) {
    return NextResponse.json({ error: 'Failed to download file' }, { status: 500 });
  }
}
