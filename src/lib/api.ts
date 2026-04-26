const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

// Helper function to get auth token
const getAuthToken = (): string | null => {
  return localStorage.getItem('token');
};

// Helper function to handle API responses
async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'An error occurred' }));
    throw new Error(error.error || `HTTP ${response.status}: ${response.statusText}`);
  }
  return response.json();
}

// Auth API
export const authAPI = {
  signup: async (email: string, password: string) => {
    const response = await fetch(`${API_BASE_URL}/auth/signup`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });
    return handleResponse(response);
  },

  login: async (email: string, password: string) => {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });
    return handleResponse<{
      access_token: string;
      refresh_token: string;
      user: { id: string; email: string };
    }>(response);
  },

  refresh: async (refreshToken: string) => {
    const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
    return handleResponse<{ access_token: string }>(response);
  },
};

// Spaces API
export interface Space {
  id: string;
  name: string;
  api_key?: string;
  created_at: string;
}

export const spacesAPI = {
  list: async (): Promise<Space[]> => {
    const token = getAuthToken();
    const response = await fetch(`${API_BASE_URL}/spaces`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    return handleResponse<Space[]>(response);
  },

  create: async (name: string): Promise<Space> => {
    const token = getAuthToken();
    const response = await fetch(`${API_BASE_URL}/spaces`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({ name }),
    });
    return handleResponse<Space>(response);
  },

  delete: async (id: string): Promise<void> => {
    const token = getAuthToken();
    const response = await fetch(`${API_BASE_URL}/spaces/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    return handleResponse(response);
  },
};

// Files API
export interface FileInfo {
  name: string;
  path: string;
  size: number;
  last_modified: string;
  is_folder: boolean;
}

export const filesAPI = {
  list: async (apiKey: string, path: string = '/'): Promise<{ files: FileInfo[]; count: number }> => {
    const response = await fetch(`${API_BASE_URL}/files?path=${encodeURIComponent(path)}`, {
      headers: {
        'X-API-Key': apiKey,
      },
    });
    return handleResponse(response);
  },

  upload: async (apiKey: string, file: File, path: string = '/'): Promise<any> => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('path', path);

    const response = await fetch(`${API_BASE_URL}/files`, {
      method: 'POST',
      headers: {
        'X-API-Key': apiKey,
      },
      body: formData,
    });
    return handleResponse(response);
  },

  download: async (apiKey: string, filePath: string): Promise<Blob> => {
    const response = await fetch(`${API_BASE_URL}/files/download?path=${encodeURIComponent(filePath)}`, {
      headers: {
        'X-API-Key': apiKey,
      },
    });
    if (!response.ok) {
      throw new Error('Failed to download file');
    }
    return response.blob();
  },

  delete: async (apiKey: string, filePath: string): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/files?path=${encodeURIComponent(filePath)}`, {
      method: 'DELETE',
      headers: {
        'X-API-Key': apiKey,
      },
    });
    return handleResponse(response);
  },

  getPresignedURL: async (apiKey: string, filePath: string): Promise<{ url: string }> => {
    const response = await fetch(`${API_BASE_URL}/files/presigned-url?path=${encodeURIComponent(filePath)}`, {
      headers: {
        'X-API-Key': apiKey,
      },
    });
    return handleResponse(response);
  },
};
