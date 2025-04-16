import React from 'react';

const BuyMeCoffee = () => {
  return (
    <div className="mt-8 py-6 px-4 rounded-lg bg-[#303842] w-full max-w-md">
      <div className="text-center">
        <h3 className="text-2xl font-bold text-white">Just want to say thanks?</h3>
        <div className="mt-4">
          <a href="/stripe/donation-checkout" className="inline-block hover:opacity-90 transform hover:-translate-y-0.5 transition-all duration-200">
            <img className="h-12" src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me a Coffee" />
          </a>
        </div>
      </div>
    </div>
  );
};

export default BuyMeCoffee;
