import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import BuyMeCoffee from '../components/BuyMeCoffee';

function HomePage() {
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const checkAuthentication = async () => {
      try {
        const response = await fetch('/api/v1/user/me');
        if (response.ok) {
          const data = await response.json();
          if (!data.user) {
            // User is not authenticated, redirect to login
            navigate('/login');
          }
        } else {
          // Error with API, assume not authenticated
          navigate('/login');
        }
      } catch (error) {
        console.error('Error checking authentication:', error);
        navigate('/login');
      } finally {
        setIsLoading(false);
      }
    };

    checkAuthentication();
  }, [navigate]);

  if (isLoading) {
    return (
      <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen">
        <div className="text-xl">Loading...</div>
      </div>
    );
  }

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen overflow-x-hidden">
      <div className="bg-[#2d3436] shadow-lg shadow-black/20 p-8 md:p-12 rounded-xl text-center max-w-md w-[90%]">
        <h1 className="text-4xl md:text-[3.5rem] font-extrabold mb-6 tracking-tight text-[#f1f2f6] leading-tight">PleFi</h1>
        
        <div className="mb-6 text-blue-400">
          <svg xmlns="http://www.w3.org/2000/svg" className="h-16 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
        </div>
        
        <p className="text-lg mb-8 text-[#f1f2f6]">You are now authenticated and have access to Plex content.</p>
        
        <div className="space-y-6">
          <button
            onClick={() => window.location.href = '/api/v1/plex/check-access'}
            className="w-full flex items-center justify-center bg-[#e5a00d] hover:bg-[#f5b82e] text-[#191a1c] font-bold py-3.5 px-8 rounded-lg shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all duration-200 text-lg"
          >
            Check Plex Access
          </button>
        </div>
        
        {/* Add the Buy Me a Coffee component */}
        <BuyMeCoffee />
      </div>
    </div>
  );
}

export default HomePage;
