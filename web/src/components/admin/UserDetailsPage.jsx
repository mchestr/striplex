import React, { useState, useEffect } from 'react';

function UserDetailsPage({ userId, onBack }) {
  const [user, setUser] = useState(null);
  const [invites, setInvites] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isRevokeModalOpen, setIsRevokeModalOpen] = useState(false);

  useEffect(() => {
    if (userId) {
      fetchUserDetails();
    }
  }, [userId]);

  const fetchUserDetails = async () => {
    setIsLoading(true);
    setError(null);
    try {
      // Fetch user details
      const userResponse = await fetch(`/api/v1/plex/users/${userId}`);
      if (!userResponse.ok) {
        throw new Error(`Error fetching user: ${userResponse.statusText}`);
      }
      
      const userData = await userResponse.json();
      setUser(userData.user);
      
      // Fetch user's invite codes
      const invitesResponse = await fetch(`/api/v1/plex/users/${userId}/invites`);
      if (invitesResponse.ok) {
        const invitesData = await invitesResponse.json();
        setInvites(invitesData.invites || []);
      }
    } catch (err) {
      setError(err.message);
      console.error('Error fetching user details:', err);
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return (
      <div className="bg-[#2d3436] rounded-lg p-6">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold">User Details</h2>
          <button
            onClick={onBack}
            className="text-gray-300 hover:text-white flex items-center"
          >
            <svg className="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
            </svg>
            Back to Users
          </button>
        </div>
        <div className="flex justify-center items-center py-10">
          <div className="loader animate-spin rounded-full h-10 w-10 border-t-2 border-b-2 border-blue-500"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-[#2d3436] rounded-lg p-6">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold">User Details</h2>
          <button
            onClick={onBack}
            className="text-gray-300 hover:text-white flex items-center"
          >
            <svg className="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
            </svg>
            Back to Users
          </button>
        </div>
        <div className="bg-red-900/20 border border-red-900 text-white p-4 rounded-lg">
          {error}
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="bg-[#2d3436] rounded-lg p-6">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold">User Details</h2>
          <button
            onClick={onBack}
            className="text-gray-300 hover:text-white flex items-center"
          >
            <svg className="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
            </svg>
            Back to Users
          </button>
        </div>
        <div className="bg-yellow-900/20 border border-yellow-900 text-white p-4 rounded-lg">
          User not found
        </div>
      </div>
    );
  }

  return (
    <div className="bg-[#2d3436] rounded-lg p-6">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">User Details</h2>
        <button
          onClick={onBack}
          className="text-gray-300 hover:text-white flex items-center"
        >
          <svg className="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
          </svg>
          Back to Users
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <div className="bg-[#3a4149] p-6 rounded-lg">
          <h3 className="text-lg font-semibold mb-4 text-white">User Information</h3>
          <div className="space-y-3">
            <div>
              <p className="text-sm text-gray-400">ID</p>
              <p className="text-white">{user.id}</p>
            </div>
            <div>
              <p className="text-sm text-gray-400">Username</p>
              <p className="text-white">{user.username}</p>
            </div>
            <div>
              <p className="text-sm text-gray-400">Email</p>
              <p className="text-white">{user.email}</p>
            </div>
            <div>
              <p className="text-sm text-gray-400">UUID</p>
              <p className="text-white">{user.uuid}</p>
            </div>
            <div>
              <p className="text-sm text-gray-400">Admin Status</p>
              <p className="text-white">
                {user.is_admin ? (
                  <span className="bg-green-900/30 text-green-400 px-2 py-1 rounded-full text-xs">Admin</span>
                ) : (
                  <span className="bg-gray-700/30 text-gray-400 px-2 py-1 rounded-full text-xs">User</span>
                )}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-400">Created At</p>
              <p className="text-white">{new Date(user.created_at).toLocaleString()}</p>
            </div>
            <div>
              <p className="text-sm text-gray-400">Updated At</p>
              <p className="text-white">{new Date(user.updated_at).toLocaleString()}</p>
            </div>
          </div>
        </div>

        <div className="bg-[#3a4149] p-6 rounded-lg">
          <h3 className="text-lg font-semibold mb-4 text-white">Active Invites</h3>
          {invites && invites.length > 0 ? (
            <div className="space-y-4">
              {invites.map((invite) => (
                <div key={invite.id} className="bg-[#2d3436] p-4 rounded-lg">
                  <div className="flex justify-between items-start">
                    <div>
                      <p className="text-white font-medium">{invite.invite_code}</p>
                      <p className="text-sm text-gray-400">Entitlement: {invite.entitlement_name}</p>
                    </div>
                    <div className="text-right">
                      <p className="text-sm text-gray-400">Used on:</p>
                      <p className="text-white">{new Date(invite.used_at).toLocaleDateString()}</p>
                    </div>
                  </div>
                  <div className="mt-2">
                    <p className="text-sm text-gray-400">
                      Expires: {invite.expires_at ? new Date(invite.expires_at).toLocaleString() : 'Never'}
                    </p>
                    <div className="mt-1">
                      {invite.expires_at ? (
                        new Date(invite.expires_at) > new Date() ? (
                          <span className="bg-green-900/30 text-green-400 px-2 py-1 rounded-full text-xs">Active</span>
                        ) : (
                          <span className="bg-red-900/30 text-red-400 px-2 py-1 rounded-full text-xs">Expired</span>
                        )
                      ) : (
                        <span className="bg-green-900/30 text-green-400 px-2 py-1 rounded-full text-xs">Permanent</span>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-400 text-center py-4">No active invites found</p>
          )}
        </div>
      </div>

      <div className="bg-[#3a4149] p-6 rounded-lg">
        <h3 className="text-lg font-semibold mb-4 text-white">Actions</h3>
        <div className="flex flex-wrap gap-3">
          <button
            className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg transition-colors"
            onClick={() => setIsRevokeModalOpen(true)}
          >
            Revoke Plex Access
          </button>
        </div>
      </div>

      {/* Revoke Access Confirmation Modal */}
      {isRevokeModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-[#2d3436] rounded-lg p-6 max-w-md w-full mx-4 shadow-xl">
            <h3 className="text-xl font-semibold mb-4 text-white">Confirm Revoke Access</h3>
            <p className="mb-6 text-gray-300">Are you sure you want to revoke this user's Plex access?</p>
            <div className="flex justify-end space-x-3">
              <button 
                onClick={() => setIsRevokeModalOpen(false)}
                className="px-4 py-2 bg-gray-600 hover:bg-gray-700 rounded-lg text-white transition-colors"
              >
                Cancel
              </button>
              <button 
                onClick={() => {
                  // This would call an API to revoke access
                  alert('This would revoke the user\'s access (API not implemented yet)');
                  setIsRevokeModalOpen(false);
                }}
                className="px-4 py-2 bg-red-600 hover:bg-red-700 rounded-lg text-white transition-colors"
              >
                Revoke Access
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default UserDetailsPage;
