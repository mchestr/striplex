import React from "react";

const WelcomeStep = ({ onNext, serverInfo }) => {
  // Dynamically generate the list of features based on available server info
  const getFeatureList = () => {
    const features = [];

    if (serverInfo.requestsUrl) {
      features.push(
        <li key="requests">How to request new movies and TV shows</li>
      );
    }

    if (serverInfo.discordServerUrl) {
      features.push(<li key="discord">How to join our community Discord</li>);
    }

    features.push(<li key="tips">Tips for the best streaming experience</li>);

    return features;
  };

  return (
    <div className="space-y-6">
      <div className="bg-green-500/20 border border-green-500/50 text-green-200 p-4 rounded-lg">
        <h2 className="text-xl font-bold mb-2">
          Welcome to {serverInfo.serverName}!
        </h2>
        <p>
          Congratulations! Your access has been granted. Let's get you set up
          with everything you need.
        </p>
      </div>

      <div className="text-left bg-[#1e272e] p-5 rounded-lg">
        <h3 className="font-bold text-xl mb-3">What You'll Learn:</h3>
        <ul className="space-y-2 list-disc list-inside">{getFeatureList()}</ul>
      </div>

      <button
        onClick={onNext}
        className="w-full py-3 bg-[#4b6bfb] hover:bg-[#3557fa] text-white font-medium rounded-lg transition-colors"
      >
        Get Started
      </button>
    </div>
  );
};

export default WelcomeStep;
