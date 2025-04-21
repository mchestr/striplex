import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';

function ClaimCodePage() {
  const [code, setCode] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState(''); // 'success' or 'error'
  const navigate = useNavigate();
  const location = useLocation();

  // Extract code from the route path (last segment of URL)
  useEffect(() => {
    // Get the path from location
    const path = location.pathname;
    // Split the path by '/' and get the last segment
    const pathSegments = path.split('/');
    const codeFromPath = pathSegments[pathSegments.length - 1];
    
    // If the last segment exists and isn't just 'claim-code', set it as the code
    if (codeFromPath && codeFromPath !== 'claim') {
      setCode(codeFromPath);
    }
  }, [location.pathname]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!code.trim()) {
      setMessage('Please enter a valid code');
      setMessageType('error');
      return;
    }

    setIsSubmitting(true);
    setMessage('');
    
    try {
      const response = await fetch('/api/v1/codes/claim', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ code: code.trim() }),
      });

      const data = await response.json();

      if (response.ok) {
        setMessage('Code claimed successfully!');
        setMessageType('success');
        // Optionally redirect after successful claim
        setTimeout(() => navigate('/'), 2000);
      } else {
        setMessage(data.error || 'Failed to claim code');
        setMessageType('error');
      }
    } catch (error) {
      console.error('Error claiming code:', error);
      setMessage('An error occurred. Please try again.');
      setMessageType('error');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen">
      <div className="bg-[#2d3436] shadow-lg shadow-black/20 p-8 md:p-12 rounded-xl text-center max-w-md w-[90%]">
        <h1 className="text-3xl font-bold mb-6">Claim Your Code</h1>
        
        {message && (
          <div className={`p-4 mb-6 rounded-lg ${
            messageType === 'success' ? 'bg-green-500/20 border border-green-500/50 text-green-200' : 
            'bg-red-500/20 border border-red-500/50 text-red-200'
          }`}>
            {message}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="mb-6">
            <label htmlFor="code" className="block text-sm font-medium text-gray-300 mb-1">
              Enter Invite Code
            </label>
            <input
              id="code"
              type="text"
              value={code}
              onChange={(e) => setCode(e.target.value)}
              className="w-full p-3 bg-[#3a4149] border border-gray-600 rounded-lg text-white focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter your code"
              disabled={isSubmitting}
            />
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className={`w-full flex items-center justify-center ${
              isSubmitting ? 'bg-[#4b6bfb]/50' : 'bg-[#4b6bfb] hover:bg-[#3557fa]'
            } text-white font-medium py-3 px-4 rounded-lg focus:outline-none transition-colors`}
          >
            {isSubmitting ? 'Claiming...' : 'Claim Code'}
          </button>
        </form>

        <div className="mt-8">
          <button 
            onClick={() => navigate('/')}
            className="text-sm text-gray-400 hover:text-gray-200 transition-colors"
          >
            Back to Home
          </button>
        </div>
      </div>
    </div>
  );
}

export default ClaimCodePage;
