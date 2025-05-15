import React from "react";

const TipsAndTricksStep = ({ onComplete, onPrev, serverInfo }) => {
  return (
    <div className="space-y-6">
      <div className="bg-blue-500/20 border border-blue-500/50 text-blue-200 p-4 rounded-lg text-left">
        <h3 className="font-bold mb-2">Pro Tips:</h3>
        <ul className="list-disc list-inside space-y-2">
          <li>Check if media already exists before requesting</li>
          <li>Download the Plex app on your devices for the best experience</li>
          <li>
            Adjust video quality to use Original quality (Maximum) whenever
            possible
          </li>
          {serverInfo.discordServerUrl ? (
            <li>For any issues, reach out on Discord!</li>
          ) : (
            ""
          )}
          <li className="text-yellow-300">
            New users have limited requests initially - use them wisely!
          </li>
        </ul>
      </div>

      <div className="bg-[#1e272e] p-4 rounded-lg text-left">
        <h3 className="font-bold text-xl mb-2">Helpful Links:</h3>
        <div className="space-y-1">
          <a
            href="https://app.plex.tv/desktop"
            className="block text-blue-400 hover:underline"
          >
            Plex Web App
          </a>
          <a
            href={serverInfo.requestsUrl}
            className="block text-blue-400 hover:underline"
          >
            {serverInfo.isLoading ? "Loading..." : "Requests Portal"}
          </a>
          <a
            href="https://support.plex.tv"
            className="block text-blue-400 hover:underline"
          >
            Plex Support
          </a>
        </div>
      </div>

      <div className="flex justify-between">
        <button
          onClick={onPrev}
          className="px-5 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
        >
          Back
        </button>
        <button
          onClick={onComplete}
          className="px-5 py-2 bg-[#e5a00d] hover:bg-[#f5b82e] text-[#191a1c] font-bold rounded-lg shadow-md hover:shadow-lg transition-all duration-200"
        >
          Complete & Go to Home
        </button>
      </div>
    </div>
  );
};

export default TipsAndTricksStep;
