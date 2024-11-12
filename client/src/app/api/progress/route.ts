import { NextResponse, NextRequest } from "next/server";

export async function GET(req: NextRequest) {
  const url = new URL(req.url);
  const id = url.searchParams.get("id");
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}//progress?id=${id}`);
  const data = await response.json();
  return NextResponse.json(data);
} 
