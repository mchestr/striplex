import React from 'react';

function Footer() {
  return (
    <div className="mt-6 mb-4 text-center text-gray-500 text-xs">
      <p className="flex justify-center items-center">
        Powered by <a href="https://github.com/mchestr/plefi" className="text-[#e5a00d] hover:text-[#f5b82e] transition-colors mx-1" target="_blank" rel="noopener noreferrer">PleFi</a> 
        &mdash;
        <a 
          href="/stripe/donation" 
          className="text-[#e5a00d] hover:text-[#f5b82e] transition-colors flex items-center ml-1"
        >
          <span>Support this project</span>
          <span className="ml-1">☕</span>
        </a>
      </p>
    </div>
  );
}

export default Footer;