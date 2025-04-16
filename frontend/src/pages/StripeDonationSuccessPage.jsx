import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

function StripeDonationSuccessPage() {
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
          ❤️
        </div>

        <h1 className="text-4xl font-extrabold mb-4">Thank You for Your Donation!</h1>
        
        <p className="text-lg mb-6 text-[#f1f2f6]">
          Your support is greatly appreciated{userData ? `, ${userData.username}` : ''}!
        </p>
        
        <p className="text-lg mb-6 text-[#f1f2f6]">
          Your generous donation helps us maintain and improve our Plex service for everyone.
        </p>
        
        <p className="text-lg mb-6 text-[#f1f2f6]">
          We <span className="inline-block animate-pulse">❤️</span> supporters like you!
        </p>

        <div className="flex flex-wrap justify-center gap-4 mt-8">
          <button 
            onClick={() => navigate('/')}
            className="px-7 py-3 bg-[#4b6bfb] hover:bg-[#3557fa] text-white font-bold rounded-lg shadow-md hover:shadow-lg transition transform hover:-translate-y-0.5"
          >
            Return Home
          </button>
        </div>
      </div>
    </div>
  );
}

export default StripeDonationSuccessPage;
