import React, { useState } from 'react';

const LoginPage = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');

  const handlePlexSignIn = async () => {
    setIsLoading(true);
    setErrorMessage('');
    
    try {
      // Calculate center position for the popup
      const width = 500;
      const height = 700;
      const left = window.screen.width / 2 - width / 2;
      const top = window.screen.height / 2 - height / 2;
      
      // Open popup window for Plex authentication
      const popup = window.open(
        '/plex/auth',
        'plexAuthWindow',
        `width=${width},height=${height},top=${top},left=${left},resizable=yes,scrollbars=yes,status=yes`
      );
      
      // Check if popup was blocked by browser
      if (!popup || popup.closed || typeof popup.closed === 'undefined') {
        setErrorMessage('Popup blocked! Please allow popups for this site and try again.');
        setIsLoading(false);
        return;
      }
      
      // Create an interval to check when the popup is closed
      const checkPopupClosed = setInterval(() => {
        if (popup.closed) {
          clearInterval(checkPopupClosed);
          setIsLoading(false);
        }
      }, 500);
      
    } catch (error) {
      setErrorMessage('Failed to connect to Plex authentication service. Please try again.');
      setIsLoading(false);
    }
  };

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen overflow-x-hidden">
      <div className="bg-[#2d3436] shadow-lg shadow-black/20 p-8 md:p-12 rounded-xl text-center max-w-md w-[90%]">
        <h1 className="text-4xl md:text-[3.5rem] font-extrabold mb-4 tracking-tight text-[#f1f2f6] leading-tight">PleFi</h1>
        
        {errorMessage && (
          <div className="mb-6 p-4 bg-red-500/20 border border-red-500/50 rounded-lg text-red-200">
            {errorMessage}
          </div>
        )}
        
        <div className="space-y-6">
          <button
            onClick={handlePlexSignIn}
            disabled={isLoading}
            className="w-full flex items-center justify-center bg-[#e5a00d] hover:bg-[#f5b82e] text-[#191a1c] font-bold py-3.5 px-8 rounded-lg shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all duration-200 text-lg disabled:opacity-70 disabled:cursor-not-allowed"
          >
            <svg 
              className="w-6 h-6 mr-2" 
              viewBox="0 0 24 24" 
              fill="currentColor" 
              xmlns="http://www.w3.org/2000/svg"
            >
              <path d="M20.9716 10.5868L4.14616 0.130286C4.14574 0.129928 4.14514 0.129928 4.14472 0.129571C3.8448 -0.0477856 3.42243 -0.0439643 3.1244 0.140786C2.8252 0.326964 2.66577 0.640857 2.66577 0.999107V23.0011C2.66577 23.3582 2.82414 23.6709 3.1224 23.8574C3.27186 23.9525 3.43606 24 3.60025 24C3.76763 24 3.93541 23.9506 4.08685 23.852L20.9716 13.4131C21.2747 13.2269 21.4444 12.9092 21.4444 12.4999C21.4444 12.0908 21.2747 11.7731 20.9716 11.5868Z" />
            </svg>
            {isLoading ? 'Connecting...' : 'Sign in with Plex'}
          </button>
        </div>
        
        {/* Buy Me a Coffee section */}
        <div className="mt-8 py-6 px-4 rounded-lg bg-[#303842] w-full max-w-md">
          <div className="text-center">
            <h3 className="text-2xl font-bold text-white">Just want to say thanks?</h3>
            <div className="mt-4">
              <a href="/stripe/donation-checkout" className="inline-block hover:opacity-90 transform hover:-translate-y-0.5 transition-all duration-200">
                <img className="h-12" src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me a Coffee" />
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
