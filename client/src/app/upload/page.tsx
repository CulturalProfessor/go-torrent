'use client';

import { useState, useEffect, ChangeEvent } from 'react';

interface DownloadHistoryItem {
  torrentID: string;
  fileName: string;
  completionDate: string;
}

export default function UploadPage() {
  const [progress, setProgress] = useState<number | null>(null);
  const [torrentID, setTorrentID] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [fileName, setFileName] = useState<string | null>(null);
  const [downloadHistory, setDownloadHistory] = useState<DownloadHistoryItem[]>([]);

  const handleUpload = async (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('torrent', file);

    setFileName(file.name);
    setIsLoading(true);
    const response = await fetch('/api/upload', { method: 'POST', body: formData });
    const data = await response.json();
    setTorrentID(data.torrentID);
    setIsLoading(false);
  };

  useEffect(() => {
    if (!torrentID) return;

    const intervalId = setInterval(async () => {
      const response = await fetch(`/api/progress?id=${torrentID}`);
      const data = await response.json();
      setProgress(data.progress);

      if (data.progress >= 100) {
        clearInterval(intervalId);
        const completionDate = new Date().toLocaleString();
        addToDownloadHistory(torrentID, fileName!, completionDate);
      }
    }, 1000);

    return () => clearInterval(intervalId);
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [torrentID]);

  const addToDownloadHistory = (id: string, name: string, date: string) => {
    setDownloadHistory((prevHistory) => [
      ...prevHistory,
      { torrentID: id, fileName: name, completionDate: date },
    ]);
  };

  const downloadFile = () => {
    if (torrentID) {
      window.open(`/api/download?id=${torrentID}`, '_blank');
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-b from-purple-600 to-indigo-800 p-6">
      <div className="max-w-lg bg-white rounded-lg shadow-lg p-8 text-center text-gray-800">
        <h1 className="text-3xl font-semibold mb-4 text-purple-700">Torrent Downloader</h1>

        <div className="mb-6">
          <a
            href="https://drive.google.com/uc?export=download&id=1At-GhxB_7PGoYAh7IsrTaEw63uiyRkJt"
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-500 hover:underline"
          >
            Download Sample Torrent File
          </a>
        </div>

        <input
          type="file"
          onChange={handleUpload}
          className="w-full py-2 px-4 mb-4 border rounded-md border-gray-300 text-gray-600"
          disabled={isLoading}
        />
        
        {isLoading && <p className="italic text-gray-500 mb-4">Uploading...</p>}
        
        {progress !== null && (
          <div className="mt-4">
            <p className="text-lg font-medium text-gray-700 mb-2">
              Download Progress: {progress.toFixed(2)}%
            </p>
            <div className="w-full bg-gray-300 rounded-full h-4 mb-2">
              <div
                className="bg-blue-600 h-4 rounded-full transition-all duration-500 ease-in-out"
                style={{ width: `${progress}%` }}
              ></div>
            </div>
          </div>
        )}
        
        {progress === 100 && (
          <button
            onClick={downloadFile}
            className="mt-4 px-6 py-2 bg-green-600 text-white font-semibold rounded-md hover:bg-green-700 transition duration-200"
          >
            Download File
          </button>
        )}
        
        {progress === null && !isLoading && (
          <p className="mt-4 text-sm text-gray-600">Please upload a torrent file to start.</p>
        )}

        {/* Download History Section */}
        {downloadHistory.length > 0 && (
          <div className="mt-8 text-left w-full">
            <h2 className="text-xl font-semibold text-purple-700 mb-2">Download History</h2>
            <ul className="space-y-2">
              {downloadHistory.map((item, index) => (
                <li key={index} className="bg-gray-100 p-3 rounded-lg shadow-sm flex flex-col">
                  <span className="font-medium text-gray-700">File: {item.fileName}</span>
                  <span className="text-gray-500 text-sm">ID: {item.torrentID}</span>
                  <span className="text-gray-500 text-sm">Completed on: {item.completionDate}</span>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}
