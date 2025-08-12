type ShortUrl = {
  _id: string;
  originalUrl: string;
  shortCode: string;
  visits: number;
  userId: string;
  lastVisited: Date;
  createdAt: Date;
};

type User = {
  _id: string;
  name: string;
  profile: string;
  email: string;
  shortUrls: Array<ShortUrl> | [];
};

type UserContextType = {
  user: User | null;
  isFetching: boolean;
  error: string;
};
