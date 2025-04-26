import React from "react";

const WelcomeStep = ({ onNext, serverInfo }) => {
  return (
    <div className="space-y-6">
      <div className="bg-green-600/10 border border-green-600/25 text-green-300/90 p-4 rounded-lg">
        <h2 className="text-lg font-medium mb-1">
          Welcome to {serverInfo.serverName}!
        </h2>
        <p className="text-sm opacity-90">
          Congratulations! Your access has been granted. Let's get you set up
          with everything you need.
        </p>
      </div>

      <div className="text-left bg-[#1e272e] p-5 rounded-lg">
        <h3 className="font-bold text-xl mb-3">
          New to Plex? Here's the Basics:
        </h3>

        <div className="mb-4">
          <h4 className="font-semibold text-lg text-green-200">
            What is Plex?
          </h4>
          <p className="ml-4">
            <a 
              href="https://www.plex.tv/about/" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-blue-400 hover:text-blue-300 underline"
            >
              Plex
            </a> is a media streaming platform that lets you access movies, TV
            shows, music, and more from any device. <a 
              href="https://www.plex.tv/watch-free/" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-blue-400 hover:text-blue-300 underline"
            >
              Learn more about what Plex offers
            </a>.
          </p>
        </div>

        <div className="mb-4">
          <h4 className="font-semibold text-lg text-green-200">
            Getting Started:
          </h4>
          <ol className="list-decimal list-inside space-y-2 ml-4">
            <li>
              Download the Plex app on your device: <a 
                href="https://www.plex.tv/download/" 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-blue-400 hover:text-blue-300 underline"
              >
                Get Plex Apps
              </a>
            </li>
            <li>
              Sign in with your account (the same one you used to register here)
            </li>
            <li>Select "{serverInfo.serverName}" from available servers</li>
            <li>
              Browse content by category, search for titles, or explore
              recommendations
            </li>
          </ol>
        </div>

        <div>
          <h4 className="font-semibold text-lg text-green-200">Pro Tips:</h4>
          <ul className="list-disc list-inside space-y-1 ml-4">
            <li>Create a watchlist for content you want to enjoy later</li>
            <li>Check for device-specific apps for the best experience</li>
            <li>Adjust video quality settings based on your internet speed</li>
            <li>
              Visit the <a 
                href="https://support.plex.tv/articles/" 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-blue-400 hover:text-blue-300 underline"
              >
                Plex Support Center
              </a> if you need help
            </li>
          </ul>
        </div>
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
