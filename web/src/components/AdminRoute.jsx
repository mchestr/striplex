import React, { useEffect, useState } from 'react';
import { Navigate, Outlet } from 'react-router-dom';

function AdminRoute() {
  const [isAdmin, setIsAdmin] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const checkAdminStatus = async () => {
      try {
        const response = await fetch('/api/v1/user/me');
        const data = await response.json();
        setIsAdmin(data.user?.is_admin || false);
      } catch (error) {
        console.error('Admin check failed:', error);
        setIsAdmin(false);
      } finally {
        setIsLoading(false);
      }
    };

    checkAdminStatus();
  }, []);

  if (isLoading) {
    return (
      <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen">
        <div className="text-xl">Checking permissions...</div>
      </div>
    );
  }

  return isAdmin ? <Outlet /> : <Navigate to="/" replace />;
}

export default AdminRoute;
