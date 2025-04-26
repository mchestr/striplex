import React, { useState } from "react";

function PlexSignInButton({ 
  onSuccess, 
  onError, 
  nextUrl = "/",
  className = "" 
}) {
  const [isLoading, setIsLoading] = useState(false);
  const [popupBlocked, setPopupBlocked] = useState(false);

  const directRedirect = () => {
    // Store current page URL to return to after authentication
    const returnTo = window.location.href;
    localStorage.setItem('plexAuthReturnTo', returnTo);
    
    // Redirect directly to Plex auth
    window.location.href = `/plex/auth?next=${nextUrl}`;
  };

  const handlePlexSignIn = async () => {
    setIsLoading(true);
    setPopupBlocked(false);

    try {
      // Calculate center position for the popup
      const width = 500;
      const height = 700;
      const left = window.screen.width / 2 - width / 2;
      const top = window.screen.height / 2 - height / 2;

      // Open popup window for Plex authentication
      const popup = window.open(
        `/plex/auth?next=/login-success`,
        "plexAuthWindow",
        `width=${width},height=${height},top=${top},left=${left},resizable=yes,scrollbars=yes,status=yes`
      );

      // Check if popup was blocked by browser
      if (!popup || popup.closed || typeof popup.closed === "undefined") {
        setIsLoading(false);
        setPopupBlocked(true);
        
        if (onError) {
          onError("Popup blocked! You can use the direct sign in option below.");
        }
        return;
      }

      // Listen for messages from the popup window
      const messageHandler = (event) => {
        if (event.data.type === "PLEX_AUTH_SUCCESS") {
          window.removeEventListener("message", messageHandler);
          setIsLoading(false);
          if (onSuccess) {
            onSuccess(event.data);
          }
        }
      };
      window.addEventListener("message", messageHandler);

      // Create an interval to check when the popup is closed
      const checkPopupClosed = setInterval(() => {
        if (popup.closed) {
          clearInterval(checkPopupClosed);
          window.removeEventListener("message", messageHandler);
          setIsLoading(false);
        }
      }, 500);
    } catch (error) {
      setIsLoading(false);
      if (onError) {
        onError("Failed to connect to Plex authentication service. Please try again.");
      }
    }
  };

  return (
    <>
      <button
        onClick={handlePlexSignIn}
        disabled={isLoading}
        className={`w-full flex items-center justify-center bg-[#e5a00d] hover:bg-[#f5b82e] text-[#191a1c] font-bold py-3.5 px-8 rounded-lg shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all duration-200 text-lg disabled:opacity-70 disabled:cursor-not-allowed ${className}`}
      >
        <svg
          className="w-6 h-6 mr-2"
          viewBox="0 0 24 24"
          fill="currentColor"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M20.9716 10.5868L4.14616 0.130286C4.14574 0.129928 4.14514 0.129928 4.14472 0.129571C3.8448 -0.0477856 3.42243 -0.0439643 3.1244 0.140786C2.8252 0.326964 2.66577 0.640857 2.66577 0.999107V23.0011C2.66577 23.3582 2.82414 23.6709 3.1224 23.8574C3.27186 23.9525 3.43606 24 3.60025 24C3.76763 24 3.93541 23.9506 4.08685 23.852L20.9716 13.4131C21.2747 13.2269 21.4444 12.9092 21.4444 12.4999C21.4444 12.0908 21.2747 11.7731 20.9716 11.5868Z" />
        </svg>
        {isLoading ? "Connecting..." : "Sign in with Plex"}
      </button>
      
      {popupBlocked && (
        <div className="mt-4">
          <p className="text-amber-300 text-sm mb-2">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 inline-block mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            Popup blocked by your browser
          </p>
          <button
            onClick={directRedirect}
            className="w-full flex items-center justify-center bg-[#2a333b] hover:bg-[#364249] text-white font-medium py-2.5 px-4 rounded-lg border border-gray-600 transition-colors"
          >
            Continue with Direct Sign In
          </button>
        </div>
      )}
    </>
  );
}

export default PlexSignInButton;
