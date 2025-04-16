import React, { useEffect } from 'react';

const LoginSuccessPage = () => {
  useEffect(() => {
    // Add a slight delay before attempting to close
    const timeoutId = setTimeout(() => {
      // This will only work if this page is opened in a popup window
      if (window.opener) {
        // Send a message to the parent window before closing
        window.opener.postMessage({ type: 'PLEX_AUTH_SUCCESS' }, '*');
      }
    }, 1500);

    return () => clearTimeout(timeoutId);
  }, []);

  const handleClose = () => {
    if (window.opener) {
      window.opener.postMessage({ type: 'PLEX_AUTH_SUCCESS' }, '*');
    }
    window.close();
  };

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen overflow-x-hidden">
      <div className="bg-[#2d3436] shadow-lg shadow-black/20 p-8 md:p-12 rounded-xl text-center max-w-md w-[90%]">
        <div className="mb-6 text-green-400">
          <svg xmlns="http://www.w3.org/2000/svg" className="h-16 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h1 className="text-3xl font-bold mb-4">Login Successful!</h1>
        <p className="text-lg mb-8">You have successfully authenticated with Plex.</p>
        <p className="text-sm text-gray-400 mb-6">You can now close this window and return to the application.</p>
        
        <button 
          onClick={handleClose}
          className="bg-green-500 hover:bg-green-600 text-white font-medium py-2 px-4 rounded-lg transition-colors duration-200"
        >
          Close Window
        </button>
      </div>
    </div>
  );
};

export default LoginSuccessPage;
