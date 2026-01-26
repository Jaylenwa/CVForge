import React from 'react';
import { Button } from './Button';

type Props = {
  children: React.ReactNode;
};

type State = {
  error: Error | null;
};

export class ErrorBoundary extends React.Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error) {
    return { error };
  }

  render() {
    if (this.state.error) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50 px-6">
          <div className="bg-white border border-gray-200 rounded-2xl shadow-sm p-6 max-w-xl w-full">
            <div className="text-lg font-bold text-gray-900">页面加载失败</div>
            <div className="mt-2 text-sm text-gray-600 break-words">{this.state.error.message || String(this.state.error)}</div>
            <div className="mt-5 flex items-center gap-3">
              <Button onClick={() => window.location.reload()}>刷新</Button>
              <Button variant="outline" onClick={() => this.setState({ error: null })}>返回</Button>
            </div>
          </div>
        </div>
      );
    }
    return this.props.children;
  }
}
