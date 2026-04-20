import { render, screen } from '@testing-library/react';
import App from './App';

test('renders uptime monitor heading', () => {
  render(<App />);
  const heading = screen.getByText(/Uptime Monitor/i);
  expect(heading).toBeInTheDocument();
});