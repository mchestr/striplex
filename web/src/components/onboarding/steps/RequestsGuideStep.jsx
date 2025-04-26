import React from "react";

const RequestsGuideStep = ({ onNext, onPrev, serverInfo, nextStepName }) => {
  return (
    <div className="space-y-6">
      <div className="bg-[#1e272e] p-6 rounded-lg text-center">
        <h3 className="font-bold text-xl mb-4">Using Requests</h3>
        <p className="mb-4">
          Can't find what you want? Request new TV shows and movies using the
          requests portal.
        </p>

        <div className="bg-amber-500/10 border border-amber-500/20 text-amber-300/80 p-3 rounded-lg text-left text-sm">
          <h3 className="font-medium mb-1">Request Limits:</h3>
          <p>
            Users are initially limited on the number of shows and movies they
            can request.
          </p>
          <p>Limits are reset weekly and will increase with more watch time.</p>
        </div>

        {serverInfo.isLoading ? (
          <div className="mt-4 text-gray-400">Loading requests link...</div>
        ) : (
          <a
            href={serverInfo.requestsUrl}
            className="inline-block mt-4 px-5 py-2 bg-[#5C7CFA] hover:bg-[#4B6BFB] text-white rounded-lg transition-colors"
            target="_blank"
            rel="noreferrer"
          >
            Open Requests Portal
          </a>
        )}
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
          className="px-5 py-2 bg-[#4b6bfb] hover:bg-[#3557fa] text-white rounded-lg transition-colors"
        >
          Next: {nextStepName}
        </button>
      </div>
    </div>
  );
};

export default RequestsGuideStep;
