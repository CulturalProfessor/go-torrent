import { NextResponse } from "next/server";

export async function GET() {
  try {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/ping`);
    if (response.ok) {
      const text = await response.text();
      return NextResponse.json({ message: text });
    } else {
      return NextResponse.json({ message: "Failed to reach the server." }, { status: 502 });
    }
  } catch (error) {
    console.error("Error connecting to the server:", error);
    return NextResponse.json({ message: "Error connecting to the server." }, { status: 500 });
  }
}
