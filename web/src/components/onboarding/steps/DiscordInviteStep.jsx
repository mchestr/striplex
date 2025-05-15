import React from "react";

const DiscordInviteStep = ({ onNext, onPrev, serverInfo, nextStepName }) => {
  return (
    <div className="space-y-6">
      <div className="bg-[#1e272e] p-6 rounded-lg text-center">
        <h3 className="font-bold text-xl mb-4">Join Our Discord Server</h3>
        <p className="mb-4">
          Get help, stay updated, and chat with fellow users. Join our Discord
          community to enhance your experience!
        </p>

        <a
          href={serverInfo.discordServerUrl}
          className="inline-block mt-4 px-5 py-2 bg-[#e5a00d] hover:bg-[#f5b82e] text-[#191a1c] font-bold rounded-lg shadow-md hover:shadow-lg transition-all duration-200"
          target="_blank"
          rel="noreferrer"
        >
          Join Discord Server
        </a>
      </div>

      <div className="flex justify-between">
        <button
          onClick={onPrev}
          className="px-5 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
        >
          Back
        </button>
        <button
          onClick={onNext}
          className="px-5 py-2 bg-[#e5a00d] hover:bg-[#f5b82e] text-[#191a1c] font-bold rounded-lg shadow-md hover:shadow-lg transition-all duration-200"
        >
          Next: {nextStepName}
        </button>
      </div>
    </div>
  );
};

export default DiscordInviteStep;
