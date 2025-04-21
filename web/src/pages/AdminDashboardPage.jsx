import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, useLocation, Routes, Route, Navigate } from 'react-router-dom';
import CodeListPage from '../components/admin/CodeListPage';
import CodeDetailsPage from '../components/admin/CodeDetailsPage';
import UserListPage from '../components/admin/UserListPage';
import UserDetailsPage from '../components/admin/UserDetailsPage';

function AdminDashboardPage() {
  const [currentSection, setCurrentSection] = useState('codes');
  const navigate = useNavigate();
  const location = useLocation();

  // Keep sidebar selection in sync with URL path
  useEffect(() => {
    if (location.pathname.includes('/admin/codes')) {
      setCurrentSection('codes');
    } else if (location.pathname.includes('/admin/users')) {
      setCurrentSection('users');
    } else if (location.pathname.includes('/admin/settings')) {
      setCurrentSection('settings');
    }
  }, [location.pathname]);

  // Handle section changes
  const handleSectionChange = (sectionId) => {
    setCurrentSection(sectionId);
    
    // Navigate to appropriate route based on section
    switch(sectionId) {
      case 'codes':
        navigate('/admin/codes');
        break;
      case 'users':
        navigate('/admin/users');
        break;
      case 'settings':
        navigate('/admin/settings');
        break;
      default:
        navigate('/admin');
    }
  };

  // Menu items for the sidebar
  const menuItems = [
    { id: 'codes', label: 'Invite Codes', icon: 'ticket' },
    { id: 'users', label: 'Users', icon: 'users' },
    { id: 'settings', label: 'Settings', icon: 'cog' }
  ];
  
  // Renders the appropriate icon based on the icon name
  const renderIcon = (iconName) => {
    switch (iconName) {
      case 'ticket':
        return (
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v10a2 2 0 002 2h14a2 2 0 002-2V7a2 2 0 00-2-2H5z"></path>
          </svg>
        );
      case 'users':
        return (
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path>
          </svg>
        );
      case 'cog':
        return (
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
          </svg>
        );
      default:
        return null;
    }
  };

  // CodeList component with navigation
  const CodeListWithNav = () => {
    const navigate = useNavigate();
    
    const handleViewCodeDetails = (codeId) => {
      navigate(`/admin/codes/${codeId}`);
    };
    
    return <CodeListPage onViewCodeDetails={handleViewCodeDetails} />;
  };

  // CodeDetails component with navigation
  const CodeDetailsWithNav = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    
    const handleBackToList = () => {
      navigate('/admin/codes');
    };
    
    return <CodeDetailsPage codeId={id} onBack={handleBackToList} />;
  };

  // UserList component with navigation
  const UserListWithNav = () => {
    const navigate = useNavigate();
    
    const handleViewUserDetails = (userId) => {
      navigate(`/admin/users/${userId}`);
    };
    
    return <UserListPage onViewUserDetails={handleViewUserDetails} />;
  };

  // UserDetails component with navigation
  const UserDetailsWithNav = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    
    const handleBackToList = () => {
      navigate('/admin/users');
    };
    
    return <UserDetailsPage userId={id} onBack={handleBackToList} />;
  };

  // Settings placeholder
  const Settings = () => (
    <div className="p-6 bg-[#2d3436] rounded-lg">
      <h2 className="text-2xl font-bold mb-4">Admin Settings</h2>
      <p className="text-gray-300">Settings functionality will be implemented here.</p>
    </div>
  );

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] min-h-screen">
      <div className="flex flex-col md:flex-row">
        {/* Sidebar */}
        <div className="bg-[#2d3436] w-full md:w-64 md:min-h-screen p-4">
          <div className="flex items-center justify-between md:justify-start mb-8">
            <h2 className="text-xl font-bold">Admin Dashboard</h2>
            <button 
              className="md:hidden text-gray-400 hover:text-white"
              onClick={() => navigate('/')}
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          </div>
          
          <nav>
            <ul className="space-y-2">
              {menuItems.map((item) => (
                <li key={item.id}>
                  <button
                    className={`w-full flex items-center p-3 rounded-lg transition-colors ${
                      currentSection === item.id 
                        ? 'bg-[#4b6bfb] text-white' 
                        : 'text-gray-300 hover:bg-[#3a4149]'
                    }`}
                    onClick={() => handleSectionChange(item.id)}
                  >
                    {renderIcon(item.icon)}
                    <span className="ml-3">{item.label}</span>
                  </button>
                </li>
              ))}
            </ul>
          </nav>
          
          <div className="mt-auto pt-8">
            <button 
              onClick={() => navigate('/')} 
              className="w-full flex items-center p-3 text-gray-300 hover:bg-[#3a4149] rounded-lg transition-colors"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
              </svg>
              <span className="ml-3">Return Home</span>
            </button>
          </div>
        </div>
        
        {/* Main Content with Routes */}
        <div className="flex-1 p-4 md:p-8 overflow-auto">
          <div className="max-w-6xl mx-auto">
            <Routes>
              <Route path="/" element={<Navigate to="/admin/codes" replace />} />
              <Route path="/codes" element={<CodeListWithNav />} />
              <Route path="/codes/:id" element={<CodeDetailsWithNav />} />
              <Route path="/users" element={<UserListWithNav />} />
              <Route path="/users/:id" element={<UserDetailsWithNav />} />
              <Route path="/settings" element={<Settings />} />
              <Route path="*" element={<Navigate to="/admin/codes" replace />} />
            </Routes>
          </div>
        </div>
      </div>
    </div>
  );
}

export default AdminDashboardPage;
