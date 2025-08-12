import { useUser } from "@/context/UserContext";
import { useEffect, useRef, useState } from "react";
import {
  Form,
  Input,
  addToast,
  Card,
  CardHeader,
  CardBody,
  Navbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
  Button,
  User,
  CardFooter,
  Divider,
  Tooltip,
} from "@heroui/react";
import axios, { type AxiosError } from "axios";
import { SERVER_URL } from "@/utils/constants";
import { useNavigate } from "react-router";
import {
  CircleNotch,
  Copy,
  ShareNetwork,
  Trash,
  ArrowBendUpRight,
  Clipboard,
  CheckCircle,
} from "phosphor-react";

export default function Dashboard() {
  const { user, isFetching, error } = useUser();
  const [longUrl, setLongUrl] = useState("");
  const [formError, setFormError] = useState({});
  const [isLoading, setIsLoading] = useState(false);
  const [shortUrl, setShortUrl] = useState("");
  const [isOpen, setIsOpen] = useState(false);
  const shortUrlRef = useRef<HTMLInputElement>(null);
  const [isCopied, setIsCopied] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    if (error) {
      addToast({
        title: "Error",
        description: error,
        color: "danger",
      });
    }
  }, []);

  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (
      !longUrl ||
      !/^(?:(?:https?):\/\/|www\.)[^\s/$.?#].[^\s]*$/.test(longUrl)
    ) {
      setFormError({ error: "Please provide a valid url" });
      return;
    }

    setIsLoading(true);
    try {
      const response = await axios.post(
        `${SERVER_URL}/urls`,
        { originalUrl: longUrl },
        { withCredentials: true }
      );

      if (response.status === 201) {
        setShortUrl(
          `${window.location.protocol}/${window.location.host}/${response.data.shortCode}`
        );
        addToast({
          title: "Creation Success",
          description: "Short url created",
          color: "success",
        });
      }
    } catch (err) {
      const error = err as AxiosError;
      addToast({
        title: "Creation Failed",
        description: (error.response?.data as any).error,
        color: "danger",
      });
    } finally {
      setIsLoading(false);
      setLongUrl("");
    }
  };

  const deleteShortUrl = async (shortCode: string) => {
    setIsLoading(true);
    try {
      const response = await axios.delete(`${SERVER_URL}/urls/${shortCode}`, {
        withCredentials: true,
      });

      if (response.status === 200) {
        addToast({
          title: "Success Delete",
          description: "Short url deleted",
          color: "success",
        });
      }
    } catch (err) {
      const error = err as AxiosError;
      addToast({
        title: "Delete Failed",
        description: (error.response?.data as any).error,
        color: "danger",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const onLogout = async () => {
    setIsLoading(true);
    try {
      const response = await axios.get(`${SERVER_URL}/auth/logout`, {
        withCredentials: true,
      });

      if (response.status === 200) {
        navigate("/login");
      }
    } catch (err) {
      const error = err as AxiosError;
      addToast({
        title: "Logout Failed",
        description: (error.response?.data as any).error,
        color: "danger",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const copyToclipboard = () => {
    if (navigator.clipboard) {
      navigator.clipboard.writeText(shortUrl);
      setIsCopied(true);
      setTimeout(() => {
        setIsCopied(false);
      }, 2000);
    } else {
      shortUrlRef.current?.select();
      document.execCommand("copy");
      setIsCopied(true);
      setTimeout(() => {
        setIsCopied(false);
      }, 2000);
    }
  };

  const handleSharing = async () => {
    try {
      await navigator.share({
        title: "Short URL",
        url: shortUrl,
      });
    } catch (error) {
      addToast({
        title: "Sharing Failed",
        description: (error as Error).message,
        color: "danger",
      });
    }
  };

  return (
    <div>
      {isFetching ? (
        <CircleNotch size={48} className="animate-spin text-orange-300" />
      ) : (
        <div>
          <Navbar>
            <NavbarBrand>
              <h2>MinURLy</h2>
            </NavbarBrand>
            <NavbarContent>
              <NavbarItem>
                <Button variant="ghost" onPress={() => setIsOpen(true)}>
                  My Urls
                </Button>
              </NavbarItem>
              <NavbarItem>
                <Button disabled={isLoading} onPress={onLogout} className="">
                  {isLoading ? (
                    <CircleNotch className="w-5 h-5 animate-spin" />
                  ) : (
                    "Log out"
                  )}
                </Button>
              </NavbarItem>
              <NavbarItem>
                <User name={user?.name} avatarProps={{ src: user?.profile }} />
              </NavbarItem>
            </NavbarContent>
          </Navbar>
          <Card>
            <CardHeader>
              <h2>Short a Long link</h2>
            </CardHeader>
            <CardBody>
              <Form validationErrors={formError} onSubmit={onSubmit}>
                <Input
                  isRequired
                  className=""
                  label="Long URL"
                  labelPlacement="outside"
                  placeholder="https://example.co.in/very-long-url"
                  type="text"
                  value={longUrl}
                  onChange={(e) => setLongUrl(e.target.value)}
                />
                <Button type="submit">Shorten URL</Button>
              </Form>
            </CardBody>
            <CardFooter>
              {shortUrl && (
                <div>
                  <Input
                    readOnly
                    type="text"
                    className=""
                    value={shortUrl}
                    endContent={
                      <Tooltip content="Copy to Clipboard" color="success">
                        <Button
                          isIconOnly
                          onPress={copyToclipboard}
                          variant="ghost"
                        >
                          {isCopied ? <CheckCircle /> : <Clipboard />}
                        </Button>
                      </Tooltip>
                    }
                  />
                </div>
              )}
            </CardFooter>
          </Card>
        </div>
      )}
      {isOpen && (
        <div>
          <h2>Short Urls</h2>
          <section>
            {user?.shortUrls.map((short) => (
              <Card key={short._id}>
                <CardHeader>
                  <h3>{`${window.location.protocol}/${window.location.host}/${short.shortCode}`}</h3>
                  <h4>{short.originalUrl}</h4>
                </CardHeader>
                <CardBody>
                  <Tooltip content="Visit URL" color="primary">
                    <Button isIconOnly>
                      <ArrowBendUpRight />
                    </Button>
                  </Tooltip>
                  <Tooltip content="Share url with others">
                    <Button isIconOnly onPress={handleSharing}>
                      <ShareNetwork />
                    </Button>
                  </Tooltip>
                  <Tooltip content="Copy to Clipboard">
                    <Button isIconOnly onPress={copyToclipboard}>
                      {isCopied ? <CheckCircle /> : <Copy />}
                    </Button>
                  </Tooltip>
                  <Tooltip content="Delete url" color="danger">
                    <Button
                      onPress={() => deleteShortUrl(short.shortCode)}
                      isIconOnly
                      variant="bordered"
                      color="danger"
                    >
                      <Trash />
                    </Button>
                  </Tooltip>
                </CardBody>
                <Divider />
                <CardFooter>
                  <p>{short.createdAt.toLocaleDateString()}</p>
                  <p>{`Visits: ${short.visits}`}</p>
                </CardFooter>
              </Card>
            ))}
          </section>
        </div>
      )}
    </div>
  );
}
