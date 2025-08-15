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
  Divider,
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
          `${window.location.protocol}//${window.location.host}/${response.data.shortCode}`
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

  const onVisit = (url: string) => {
    window.open(
      `${window.location.protocol}//${window.location.host}/${url}`,
      "_blank",
      "noopener,noreferrer"
    );
  };

  return (
    <>
      {isFetching ? (
        <div className="w-full h-screen flex justify-center items-center">
          <CircleNotch size={48} className="animate-spin text-orange-300" />
        </div>
      ) : (
        <div className="bg-gradient-to-tr from-orange-400 to-red-300 w-full h-screen flex flex-col items-center relative">
          <Navbar maxWidth="xl" className="flex items-center">
            <NavbarBrand className="flex-1">
              <h2 className="text-2xl font-semibold text-green-500">MinURLy</h2>
            </NavbarBrand>
            <NavbarContent justify="end">
              <NavbarItem>
                <Button
                  variant="light"
                  className="text-lg font-medium text-slate-600"
                  onPress={() => setIsOpen(!isOpen)}
                >
                  My Urls
                </Button>
              </NavbarItem>
              <NavbarItem>
                <Button
                  disabled={isLoading}
                  isLoading={isLoading}
                  spinner={<CircleNotch />}
                  color="danger"
                  radius="sm"
                  onPress={onLogout}
                  className="text-lg font-medium"
                >
                  {isLoading ? "" : "Log out"}
                </Button>
              </NavbarItem>
              <NavbarItem>
                <User name={user?.name} avatarProps={{ src: user?.profile }} />
              </NavbarItem>
            </NavbarContent>
          </Navbar>
          <div className="w-full h-full flex justify-center items-center">
            <Card className="w-1/2 px-10 py-6">
              <CardHeader>
                <h2 className="text-3xl font-semibold text-slate-700">
                  Short a Long Link
                </h2>
              </CardHeader>
              <CardBody>
                <Form
                  validationErrors={formError}
                  onSubmit={onSubmit}
                  className="flex flex-col items-start space-y-1"
                >
                  <label
                    htmlFor="longurl"
                    className="text-lg font-medium text-slate-600"
                  >
                    Long URL<i className="text-red-500">*</i>
                  </label>
                  <Input
                    isRequired
                    size="lg"
                    id="longurl"
                    variant="bordered"
                    radius="sm"
                    className="caret-orange-500 text-base font-medium text-slate-600"
                    labelPlacement="outside"
                    placeholder="https://example.co.in/very-long-url"
                    type="text"
                    value={longUrl}
                    onChange={(e) => setLongUrl(e.target.value)}
                  />
                  <Button
                    type="submit"
                    isLoading={isLoading}
                    color="success"
                    radius="sm"
                    className="mt-4 text-white font-semibold text-lg"
                  >
                    Shorten URL
                  </Button>
                </Form>
              </CardBody>
              <CardFooter className="w-full">
                {shortUrl && (
                  <div className="w-full">
                    <Input
                      ref={shortUrlRef}
                      readOnly
                      size="lg"
                      type="text"
                      radius="sm"
                      className="w-full"
                      value={shortUrl}
                      endContent={
                        <Tooltip content="Copy to Clipboard" color="success">
                          <Button
                            isIconOnly
                            onPress={copyToclipboard}
                            variant="light"
                          >
                            {isCopied ? (
                              <CheckCircle
                                size={24}
                                className="text-green-500"
                              />
                            ) : (
                              <Clipboard size={24} />
                            )}
                          </Button>
                        </Tooltip>
                      }
                    />
                  </div>
                )}
              </CardFooter>
            </Card>
          </div>
          <footer className="w-full flex px-10 justify-start items-center bg-gray-300">
            <p className="text-sm text-slate-600">
              Copyright@MinURLy {new Date().getFullYear()}
            </p>
          </footer>
        </div>
      )}

      <div
        className={`fixed bg-white z-50 h-screen shadow-xl top-0 right-0 px-4 py-6 transition-transform delay-300 ${
          isOpen ? "translate-x-0" : "translate-x-[400px]"
        }`}
      >
        <h2 className="text-3xl font-semibold text-slate-700">Short Urls</h2>
        <Divider />
        <section className="flex flex-col items-start space-y-4 mt-10">
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
    </>
  );
}
