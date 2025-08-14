import type { ShortUrl } from "@/types/types";
import { SERVER_URL } from "@/utils/constants";
import {
  Card,
  CardHeader,
  CardBody,
  CardFooter,
  Button,
  Tooltip,
  Divider,
  addToast,
} from "@heroui/react";
import axios, { AxiosError } from "axios";
import {
  Copy,
  ShareNetwork,
  Trash,
  ArrowBendUpRight,
  CheckCircle,
} from "phosphor-react";

type URLCardProps = {
  short: ShortUrl;
  isCopied: boolean;
  onVisit: () => void;
  handleSharing: () => void;
  copyToclipboard: () => void;
  deleteUrlCard: (shortCode: string) => void;
};

export function URLCard({
  short,
  isCopied,
  onVisit,
  handleSharing,
  copyToclipboard,
  deleteUrlCard,
}: URLCardProps) {
  const deleteShortUrl = async () => {
    try {
      const response = await axios.delete(
        `${SERVER_URL}/urls/${short.shortCode}`,
        {
          withCredentials: true,
        }
      );

      if (response.status === 200) {
        addToast({
          title: "Success Delete",
          description: "Short url deleted",
          color: "success",
        });

        deleteUrlCard(short.shortCode);
      }
    } catch (err) {
      const error = err as AxiosError;
      addToast({
        title: "Delete Failed",
        description: (error.response?.data as any).error,
        color: "danger",
      });
    }
  };

  return (
    <Card>
      <CardHeader>
        <h3>{`${window.location.protocol}/${window.location.host}/${short.shortCode}`}</h3>
        <h4>{short.originalUrl}</h4>
      </CardHeader>
      <CardBody>
        <Tooltip content="Visit URL" color="primary">
          <Button isIconOnly onPress={onVisit}>
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
            onPress={deleteShortUrl}
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
        <p>{new Date(short.createdAt).toLocaleDateString()}</p>
        <p>{`Visits: ${short.visits}`}</p>
      </CardFooter>
    </Card>
  );
}
