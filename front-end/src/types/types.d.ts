type ShortUrl = {
  _id: string;
  originalUrl: string;
  shortCode: string;
  visits: number;
  userId: string;
  lastVisited: string;
  createdAt: string;
};

type User = {
  id: string;
  name: string;
  profile: string;
  email: string;
  shorturls: Array<ShortUrl> | [];
};

type UserContextType = {
  user: User | null;
  isFetching: boolean;
  error: string;
};
