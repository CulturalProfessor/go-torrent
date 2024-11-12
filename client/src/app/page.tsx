'use client';

import { useEffect, useState } from "react";

export default function Home() {
  const [serverStatus, setServerStatus] = useState<string | null>(null);

  const checkServerStatus = async () => {
    try {
      const response = await fetch("/api/ping");
      if (response.ok) {
        
        setServerStatus("Server is up and running");
      } else {
        setServerStatus("Server is not reachable");
      }
    } catch (error) {
      setServerStatus("Error connecting to server");
      console.error("Ping error:", error);
    }
  };

  useEffect(() => {
    checkServerStatus();
  }, []);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-b from-blue-600 to-blue-800 text-white p-6">
      <main className="bg-white text-gray-900 p-8 rounded-lg shadow-lg text-center max-w-md w-full">
        <h1 className="text-4xl font-bold mb-4">Go Torrent Project</h1>
        <p className="text-lg mb-6">
          Upload torrents, monitor download progress, and retrieve completed files effortlessly.
        </p>

        <div className="bg-gray-100 p-4 rounded-md mb-4">
          <h2 className="text-xl font-semibold">How It Works:</h2>
          <ol className="list-decimal list-inside text-left space-y-2 mt-2 text-gray-700">
            <li>Upload a torrent file using the interface.</li>
            <li>Monitor the real-time download progress.</li>
            <li>Download the completed file.</li>
          </ol>
        </div>

        {serverStatus && (
          <div className="bg-blue-100 text-blue-800 p-3 rounded-md mt-4">
            <p>{serverStatus}</p>
          </div>
        )}

        <a
          href="/upload"
          className="mt-8 inline-block px-6 py-2 bg-blue-600 text-white rounded-full font-semibold hover:bg-blue-700 transition"
        >
          Get Started
        </a>
      </main>
    </div>
  );
}
