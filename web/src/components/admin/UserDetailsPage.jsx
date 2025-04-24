import React, { useState, useEffect } from "react";

function UserDetailsPage({ userId, onBack }) {
  const [user, setUser] = useState(null);
  const [invites, setInvites] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isRevokeModalOpen, setIsRevokeModalOpen] = useState(false);
  const [isRevoking, setIsRevoking] = useState(false);
  const [revokeError, setRevokeError] = useState(null);
  const [plexAccess, setPlexAccess] = useState(null);

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

      // Fetch user's Plex access status
      const accessResponse = await fetch(`/api/v1/plex/users/${userId}/access`);
      if (accessResponse.ok) {
        const accessData = await accessResponse.json();
        setPlexAccess(accessData.has_access || false);
      }

      // Fetch user's invite codes
      const invitesResponse = await fetch(
        `/api/v1/plex/users/${userId}/invites`
      );
      if (invitesResponse.ok) {
        const invitesData = await invitesResponse.json();
        setInvites(invitesData.invites || []);
      }
    } catch (err) {
      setError(err.message);
      console.error("Error fetching user details:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleRevokeAccess = async () => {
    setIsRevoking(true);
    setRevokeError(null);
    try {
      const response = await fetch(`/api/v1/plex/users/${userId}/access`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Failed to revoke access");
      }

      // If successful, close the modal and refresh user details
      setIsRevokeModalOpen(false);
      fetchUserDetails();
    } catch (err) {
      setRevokeError(err.message);
    } finally {
      setIsRevoking(false);
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
            <svg
              className="w-5 h-5 mr-1"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M10 19l-7-7m0 0l7-7m-7 7h18"
              ></path>
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
            <svg
              className="w-5 h-5 mr-1"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M10 19l-7-7m0 0l7-7m-7 7h18"
              ></path>
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
            <svg
              className="w-5 h-5 mr-1"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M10 19l-7-7m0 0l7-7m-7 7h18"
              ></path>
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
          <svg
            className="w-5 h-5 mr-1"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M10 19l-7-7m0 0l7-7m-7 7h18"
            ></path>
          </svg>
          Back to Users
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <div className="bg-[#3a4149] p-6 rounded-lg">
          <h3 className="text-lg font-semibold mb-4 text-white flex items-center">
            <svg
              className="w-5 h-5 mr-2"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
              ></path>
            </svg>
            User Information
          </h3>

          {/* Status indicators row */}
          <div className="flex flex-wrap gap-2 mb-4">
            <div
              className={`px-3 py-1 rounded-md text-sm font-medium ${
                user.is_admin
                  ? "bg-indigo-900/40 text-indigo-300 border border-indigo-700"
                  : "bg-gray-700/30 text-gray-400 border border-gray-600"
              }`}
            >
              {user.is_admin ? "Admin Account" : "Standard User"}
            </div>
            <div
              className={`px-3 py-1 rounded-md text-sm font-medium ${
                plexAccess
                  ? "bg-green-900/30 text-green-300 border border-green-700"
                  : "bg-red-900/30 text-red-300 border border-red-700"
              }`}
            >
              {plexAccess ? "Plex Access: Active" : "Plex Access: None"}
            </div>
          </div>

          {/* User information grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-x-4 gap-y-3">
            {/* First column */}
            <div className="space-y-3">
              <div className="bg-[#2d3436] p-3 rounded-md">
                <p className="text-xs uppercase tracking-wider text-blue-400 font-semibold mb-1">
                  Username
                </p>
                <p className="text-white text-lg font-medium">
                  {user.username}
                </p>
              </div>

              <div className="bg-[#2d3436] p-3 rounded-md">
                <p className="text-xs uppercase tracking-wider text-blue-400 font-semibold mb-1">
                  Email
                </p>
                <p className="text-white break-all">{user.email}</p>
              </div>

              <div className="bg-[#2d3436] p-3 rounded-md">
                <p className="text-xs uppercase tracking-wider text-blue-400 font-semibold mb-1">
                  Account ID
                </p>
                <p className="text-gray-300 font-mono text-sm">{user.id}</p>
              </div>
            </div>

            {/* Second column */}
            <div className="space-y-3">
              <div className="bg-[#2d3436] p-3 rounded-md">
                <p className="text-xs uppercase tracking-wider text-blue-400 font-semibold mb-1">
                  UUID
                </p>
                <p className="text-gray-300 font-mono text-xs break-all">
                  {user.uuid}
                </p>
              </div>

              <div className="bg-[#2d3436] p-3 rounded-md">
                <div className="flex justify-between">
                  <div>
                    <p className="text-xs uppercase tracking-wider text-blue-400 font-semibold mb-1">
                      Created
                    </p>
                    <p className="text-gray-300 text-sm">
                      {new Date(user.created_at).toLocaleDateString()}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="text-xs uppercase tracking-wider text-blue-400 font-semibold mb-1">
                      Updated
                    </p>
                    <p className="text-gray-300 text-sm">
                      {new Date(user.updated_at).toLocaleDateString()}
                    </p>
                  </div>
                </div>
              </div>

              <div className="bg-[#2d3436] p-3 rounded-md">
                <p className="text-xs uppercase tracking-wider text-blue-400 font-semibold mb-1">
                  Account Age
                </p>
                <p className="text-gray-300 text-sm">
                  {Math.floor(
                    (new Date() - new Date(user.created_at)) /
                      (1000 * 60 * 60 * 24)
                  )}{" "}
                  days
                </p>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-[#3a4149] p-6 rounded-lg">
          <h3 className="text-lg font-semibold mb-4 text-white">
            Active Invites
          </h3>
          {invites && invites.length > 0 ? (
            <div className="space-y-4">
              {invites.map((invite) => (
                <div key={invite.id} className="bg-[#2d3436] p-4 rounded-lg">
                  <div className="flex justify-between items-start">
                    <div>
                      <p className="text-white font-medium">
                        <a 
                          href={`/admin/codes/${invite.id}`} 
                          className="text-blue-400 hover:text-blue-300 hover:underline font-medium flex items-center"
                        >
                          {invite.invite_code}
                          <svg className="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path>
                          </svg>
                        </a>
                      </p>
                      <p className="text-sm text-gray-400">
                        Entitlement: {invite.entitlement_name}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="text-sm text-gray-400">Used on:</p>
                      <p className="text-white">
                        {new Date(invite.used_at).toLocaleDateString()}
                      </p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-400 text-center py-4">
              No active invites found
            </p>
          )}
        </div>
      </div>
      {!user.is_admin && plexAccess && (
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
      )}

      {/* Revoke Access Confirmation Modal */}
      {isRevokeModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-[#2d3436] rounded-lg p-6 max-w-md w-full mx-4 shadow-xl">
            <h3 className="text-xl font-semibold mb-4 text-white">
              Confirm Revoke Access
            </h3>
            <p className="mb-6 text-gray-300">
              Are you sure you want to revoke this user's Plex access?
            </p>
            {revokeError && (
              <div className="mb-4 p-3 bg-red-900/20 border border-red-900 text-red-400 rounded-lg">
                {revokeError}
              </div>
            )}
            <div className="flex justify-end space-x-3">
              <button
                onClick={() => setIsRevokeModalOpen(false)}
                className="px-4 py-2 bg-gray-600 hover:bg-gray-700 rounded-lg text-white transition-colors"
                disabled={isRevoking}
              >
                Cancel
              </button>
              <button
                onClick={handleRevokeAccess}
                className="px-4 py-2 bg-red-600 hover:bg-red-700 rounded-lg text-white transition-colors flex items-center"
                disabled={isRevoking}
              >
                {isRevoking ? (
                  <>
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"></div>
                    Revoking...
                  </>
                ) : (
                  "Revoke Access"
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default UserDetailsPage;
