// Mock Authentication Service

export const sendVerificationCode = async (email: string): Promise<boolean> => {
  // Simulate network delay
  await new Promise(resolve => setTimeout(resolve, 1500));
  
  // In a real app, this would call an API endpoint
  console.log(`Verification code sent to ${email}: 123456`);
  return true;
};

export const verifyCode = async (email: string, code: string): Promise<boolean> => {
  await new Promise(resolve => setTimeout(resolve, 1000));
  // Mock validation: "123456" is the magic code
  return code === '123456';
};

export const loginUser = async (email: string, password: string): Promise<{ success: boolean; token?: string }> => {
  await new Promise(resolve => setTimeout(resolve, 1000));
  // Accept any valid email format and password length > 5 for demo
  if (email.includes('@') && password.length >= 6) {
    return { success: true, token: 'mock-jwt-token' };
  }
  return { success: false };
};
