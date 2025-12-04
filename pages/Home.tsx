import React from 'react';
import { Link } from 'react-router-dom';
import { ArrowRight, CheckCircle, Cpu, FileCheck, Share2 } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { AppRoute } from '../types';

export const Home: React.FC = () => {
  return (
    <div className="flex flex-col">
      {/* Hero Section */}
      <section className="relative bg-white overflow-hidden">
        <div className="max-w-7xl mx-auto">
          <div className="relative z-10 pb-8 bg-white sm:pb-16 md:pb-20 lg:max-w-2xl lg:w-full lg:pb-28 xl:pb-32 pt-20 px-4 sm:px-6 lg:px-8">
            <main className="mt-10 mx-auto max-w-7xl sm:mt-12 md:mt-16 lg:mt-20 xl:mt-28">
              <div className="sm:text-center lg:text-left">
                <h1 className="text-4xl tracking-tight font-extrabold text-gray-900 sm:text-5xl md:text-6xl">
                  <span className="block xl:inline">Craft your perfect resume</span>{' '}
                  <span className="block text-blue-600 xl:inline">with AI precision</span>
                </h1>
                <p className="mt-3 text-base text-gray-500 sm:mt-5 sm:text-lg sm:max-w-xl sm:mx-auto md:mt-5 md:text-xl lg:mx-auto">
                  Build professional, ATS-friendly resumes in minutes. Choose from our curated templates, use AI to polish your text, and land your dream job.
                </p>
                <div className="mt-5 sm:mt-8 sm:flex sm:justify-center lg:justify-start">
                  <div className="rounded-md shadow">
                    <Link to={AppRoute.Templates}>
                      <Button size="lg" className="w-full">
                        Create Resume <ArrowRight className="ml-2 h-5 w-5" />
                      </Button>
                    </Link>
                  </div>
                  <div className="mt-3 sm:mt-0 sm:ml-3">
                     <Link to={AppRoute.Dashboard}>
                        <Button size="lg" variant="outline" className="w-full">
                            My Resumes
                        </Button>
                    </Link>
                  </div>
                </div>
              </div>
            </main>
          </div>
        </div>
        <div className="lg:absolute lg:inset-y-0 lg:right-0 lg:w-1/2">
          <img
            className="h-56 w-full object-cover sm:h-72 md:h-96 lg:w-full lg:h-full opacity-90"
            src="https://images.unsplash.com/photo-1586281380349-632531db7ed4?ixlib=rb-1.2.1&auto=format&fit=crop&w=1950&q=80"
            alt="Resume planning"
          />
        </div>
      </section>

      {/* Features Section */}
      <section className="py-12 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="lg:text-center">
            <h2 className="text-base text-blue-600 font-semibold tracking-wide uppercase">Features</h2>
            <p className="mt-2 text-3xl leading-8 font-extrabold tracking-tight text-gray-900 sm:text-4xl">
              Everything you need to get hired
            </p>
          </div>

          <div className="mt-10">
            <dl className="space-y-10 md:space-y-0 md:grid md:grid-cols-2 md:gap-x-8 md:gap-y-10">
              <div className="relative">
                <dt>
                  <div className="absolute flex items-center justify-center h-12 w-12 rounded-md bg-blue-500 text-white">
                    <Cpu size={24} />
                  </div>
                  <p className="ml-16 text-lg leading-6 font-medium text-gray-900">AI Writer & Polisher</p>
                </dt>
                <dd className="mt-2 ml-16 text-base text-gray-500">
                  Stuck on words? Use our Gemini-powered AI to generate summaries and improve your bullet points instantly.
                </dd>
              </div>

              <div className="relative">
                <dt>
                  <div className="absolute flex items-center justify-center h-12 w-12 rounded-md bg-blue-500 text-white">
                    <FileCheck size={24} />
                  </div>
                  <p className="ml-16 text-lg leading-6 font-medium text-gray-900">ATS Friendly Templates</p>
                </dt>
                <dd className="mt-2 ml-16 text-base text-gray-500">
                  Our templates are designed to pass through Applicant Tracking Systems, ensuring your resume reaches human eyes.
                </dd>
              </div>

              <div className="relative">
                <dt>
                  <div className="absolute flex items-center justify-center h-12 w-12 rounded-md bg-blue-500 text-white">
                    <Share2 size={24} />
                  </div>
                  <p className="ml-16 text-lg leading-6 font-medium text-gray-900">Share Online</p>
                </dt>
                <dd className="mt-2 ml-16 text-base text-gray-500">
                  Publish your resume with a unique link. Track views and share directly to LinkedIn or via QR code.
                </dd>
              </div>

              <div className="relative">
                <dt>
                  <div className="absolute flex items-center justify-center h-12 w-12 rounded-md bg-blue-500 text-white">
                    <CheckCircle size={24} />
                  </div>
                  <p className="ml-16 text-lg leading-6 font-medium text-gray-900">Real-time Preview</p>
                </dt>
                <dd className="mt-2 ml-16 text-base text-gray-500">
                  See changes instantly as you type. Drag and drop sections to reorganize your story effortlessly.
                </dd>
              </div>
            </dl>
          </div>
        </div>
      </section>
    </div>
  );
};
