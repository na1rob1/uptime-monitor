import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
import App from './App';

const mockSites = [
  {
    id: 1,
    url: 'https://example.com',
    name: 'Example',
    status: true,
    uptime: 99.5,
    checked_at: '2026-06-08T10:00:00Z',
    created_at: '2026-06-01T10:00:00Z',
  },
  {
    id: 2,
    url: 'https://down.com',
    name: 'Down',
    status: false,
    uptime: null,
    checked_at: null,
    created_at: '2026-06-02T10:00:00Z',
  },
];

beforeEach(() => {
  global.fetch = jest.fn((url, opts) => {
    if (!opts || opts.method === 'GET' || !opts.method) {
      return Promise.resolve({ json: () => Promise.resolve(mockSites) });
    }
    return Promise.resolve({ json: () => Promise.resolve({}) });
  });
});

afterEach(() => {
  jest.resetAllMocks();
});

test('renders uptime monitor heading', async () => {
  await act(async () => { render(<App />); });
  expect(screen.getByText(/Uptime Monitor/i)).toBeInTheDocument();
});

test('loads and displays sites', async () => {
  await act(async () => { render(<App />); });
  await waitFor(() => expect(screen.getByText('Example')).toBeInTheDocument());
  expect(screen.getByText('Down')).toBeInTheDocument();
  expect(screen.getByText('UP')).toBeInTheDocument();
  expect(screen.getByText('DOWN')).toBeInTheDocument();
  expect(screen.getByText('99.50%')).toBeInTheDocument();
});

test('handles null sites response', async () => {
  global.fetch = jest.fn(() => Promise.resolve({ json: () => Promise.resolve(null) }));
  await act(async () => { render(<App />); });
  expect(screen.getByText(/Uptime Monitor/i)).toBeInTheDocument();
});

test('add site triggers POST', async () => {
  await act(async () => { render(<App />); });
  await waitFor(() => screen.getByText('Example'));

  const urlInput = screen.getByPlaceholderText('URL');
  const nameInput = screen.getByPlaceholderText('Name');
  const addBtn = screen.getByText('Add');

  fireEvent.change(urlInput, { target: { value: 'https://new.com' } });
  fireEvent.change(nameInput, { target: { value: 'New' } });

  await act(async () => { fireEvent.click(addBtn); });

  const postCall = global.fetch.mock.calls.find(c => c[1]?.method === 'POST');
  expect(postCall).toBeTruthy();
  expect(JSON.parse(postCall[1].body)).toEqual({ url: 'https://new.com', name: 'New' });
});

test('delete site triggers DELETE', async () => {
  await act(async () => { render(<App />); });
  await waitFor(() => screen.getByText('Example'));

  const deleteBtns = screen.getAllByText('Delete');
  await act(async () => { fireEvent.click(deleteBtns[0]); });

  const deleteCall = global.fetch.mock.calls.find(c => c[1]?.method === 'DELETE');
  expect(deleteCall).toBeTruthy();
  expect(deleteCall[0]).toContain('?id=1');
});

test('edit fills inputs and disables URL', async () => {
  await act(async () => { render(<App />); });
  await waitFor(() => screen.getByText('Example'));

  await act(async () => { fireEvent.click(screen.getAllByText('Edit')[0]); });

  const urlInput = screen.getByPlaceholderText('URL');
  const nameInput = screen.getByPlaceholderText('Name');
  expect(urlInput.value).toBe('https://example.com');
  expect(nameInput.value).toBe('Example');
  expect(urlInput.disabled).toBe(true);
  expect(screen.getByText('Save')).toBeInTheDocument();
  expect(screen.getByText('Cancel')).toBeInTheDocument();
});

test('save edit triggers PUT', async () => {
  await act(async () => { render(<App />); });
  await waitFor(() => screen.getByText('Example'));

  fireEvent.click(screen.getAllByText('Edit')[0]);

  const nameInput = screen.getByPlaceholderText('Name');
  fireEvent.change(nameInput, { target: { value: 'Renamed' } });

  await act(async () => { fireEvent.click(screen.getByText('Save')); });

  const putCall = global.fetch.mock.calls.find(c => c[1]?.method === 'PUT');
  expect(putCall).toBeTruthy();
  expect(JSON.parse(putCall[1].body)).toEqual({ id: 1, url: 'https://example.com', name: 'Renamed' });
});

test('cancel edit clears inputs', async () => {
  await act(async () => { render(<App />); });
  await waitFor(() => screen.getByText('Example'));

  await act(async () => { fireEvent.click(screen.getAllByText('Edit')[0]); });
  expect(screen.getByPlaceholderText('URL').value).toBe('https://example.com');

  await act(async () => { fireEvent.click(screen.getByText('Cancel')); });
  expect(screen.getByPlaceholderText('URL').value).toBe('');
  expect(screen.getByPlaceholderText('Name').value).toBe('');
  expect(screen.getByText('Add')).toBeInTheDocument();
});

test('displays dash for null uptime and checked_at', async () => {
  await act(async () => { render(<App />); });
  await waitFor(() => screen.getByText('Down'));
  const dashes = screen.getAllByText('—');
  expect(dashes.length).toBeGreaterThan(0);
});