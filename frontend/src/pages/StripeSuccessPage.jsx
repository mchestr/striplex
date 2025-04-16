import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

function StripeSuccessPage() {
  const [userData, setUserData] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    // Fetch user data
    const fetchUserData = async () => {
      try {
        const response = await fetch('/api/v1/user/me');
        if (response.ok) {
          const data = await response.json();
          setUserData(data.user);
        }
      } catch (error) {
        console.error('Error fetching user data:', error);
      }
    };

    fetchUserData();
  }, []);

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] min-h-screen py-8 px-4">
      <div className="max-w-3xl mx-auto text-center p-12 rounded-xl shadow-lg bg-[#2d3436] shadow-black/20">
        <div className="flex items-center justify-center w-20 h-20 mx-auto mb-6 rounded-full text-4xl font-light bg-[#2b8a3e] text-[#e3f9e5]">
          ✓
        </div>

        <h1 className="text-4xl font-extrabold mb-4">Subscription Successful!</h1>
        
        <p className="text-lg mb-6 text-[#f1f2f6]">
          Thank you{userData ? `, ${userData.username}` : ''}! Your subscription will be activated soon.
        </p>
        
        <p className="text-lg mb-6 text-[#f1f2f6]">
          An invite will be sent to your Plex account shortly.
        </p>
        
        <p className="text-lg mb-6 text-[#f1f2f6]">
          You may need to accept the invite within your Plex account to gain full access.
        </p>

        <div className="flex flex-wrap justify-center gap-4 mt-8">
          <button 
            onClick={() => navigate('/')}
            className="px-7 py-3 bg-[#4b6bfb] hover:bg-[#3557fa] text-white font-bold rounded-lg shadow-md hover:shadow-lg transition transform hover:-translate-y-0.5"
          >
            Return Home
          </button>
          
          <a 
            href="https://app.plex.tv/desktop/#!/settings/manage-library-access"
            target="_blank"
            rel="noopener noreferrer"
            className="px-7 py-3 bg-[#ffb142] hover:bg-[#ff9f1a] text-white font-bold rounded-lg shadow-md hover:shadow-lg transition transform hover:-translate-y-0.5"
          >
            Check Plex Requests
          </a>
        </div>
      </div>
    </div>
  );
}

export default StripeSuccessPage;
