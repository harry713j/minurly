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
  Tooltip,
} from "@heroui/react";
import axios, { type AxiosError } from "axios";
import { SERVER_URL } from "@/utils/constants";
import { useNavigate } from "react-router";
import { CircleNotch, Clipboard, CheckCircle } from "phosphor-react";
import { URLCard } from "@/components";

export default function Dashboard() {
  const { user, isFetching, error, setUser, refetchUser } = useUser();
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

        await refetchUser();

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

  const deleteUrlCard = (shortCode: string) => {
    if (!user) return;

    setUser((prev) =>
      prev
        ? {
            ...prev,
            shorturls: prev.shorturls.filter(
              (url) => url.shortCode !== shortCode
            ),
          }
        : null
    );
  };

  const onLogout = async () => {
    setIsLoading(true);
    try {
      const response = await axios.get(`${SERVER_URL}/auth/google/logout`, {
        withCredentials: true,
      });

      if (response.status === 200) {
        setUser(null);
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
      addToast({
        title: "URL copied Successfully",
        color: "success",
      });
      setTimeout(() => {
        setIsCopied(false);
      }, 2000);
    } else {
      shortUrlRef.current?.select();
      document.execCommand("copy");
      setIsCopied(true);
      addToast({
        title: "URL copied Successfully",
        color: "success",
      });
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

  const onVisit = () => {
    window.open(`${shortUrl}`, "_blank", "noopener,noreferrer");
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
                <Button variant="light" onPress={() => setIsOpen(!isOpen)}>
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
                    ref={shortUrlRef}
                    readOnly
                    type="text"
                    className=""
                    value={shortUrl}
                    endContent={
                      <Tooltip content="Copy to Clipboard" color="success">
                        <Button
                          isIconOnly
                          onPress={copyToclipboard}
                          variant="light"
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
            {user?.shorturls?.map((short) => (
              <URLCard
                key={short._id}
                short={short}
                isCopied={isCopied}
                copyToclipboard={copyToclipboard}
                onVisit={onVisit}
                handleSharing={handleSharing}
                deleteUrlCard={deleteUrlCard}
              />
            ))}
          </section>
        </div>
      )}
    </div>
  );
}
