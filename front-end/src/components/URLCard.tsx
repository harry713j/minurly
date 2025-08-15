import type { ShortUrl } from "@/types/types";
import { SERVER_URL } from "@/utils/constants";
import {
  Card,
  CardHeader,
  CardBody,
  CardFooter,
  Button,
  Tooltip,
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
  onVisit: (url: string) => void;
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
    <Card className="flex flex-col items-start px-4 py-2 space-y-2">
      <CardHeader className="flex flex-col items-start space-y-1">
        <h3 className="font-medium text-slate-600">{`${window.location.protocol}/${window.location.host}/${short.shortCode}`}</h3>
        <h4 className="text-sm text-slate-400">{short.originalUrl}</h4>
      </CardHeader>
      <CardBody className="w-full flex flex-row items-center justify-center gap-2">
        <Tooltip content="Visit URL" color="primary">
          <Button
            isIconOnly
            onPress={() => onVisit(short.shortCode)}
            variant="flat"
            color="primary"
            className=""
          >
            <ArrowBendUpRight size={24} />
          </Button>
        </Tooltip>
        <Tooltip content="Share url with others" color="warning">
          <Button isIconOnly onPress={handleSharing} className="bg-yellow-300">
            <ShareNetwork size={24} className="text-slate-500" />
          </Button>
        </Tooltip>
        <Tooltip content="Copy to Clipboard" color="success">
          <Button
            isIconOnly
            onPress={copyToclipboard}
            color="success"
            variant="bordered"
          >
            {isCopied ? <CheckCircle size={24} /> : <Copy size={24} />}
          </Button>
        </Tooltip>
        <Tooltip content="Delete url" color="danger">
          <Button
            onPress={deleteShortUrl}
            isIconOnly
            variant="bordered"
            color="danger"
          >
            <Trash size={24} />
          </Button>
        </Tooltip>
      </CardBody>
      <CardFooter className="flex items-center justify-around border-t border-slate-300">
        <p className="text-xs text-slate-500">
          {new Date(short.createdAt).toLocaleDateString()}
        </p>
        <p className="text-xs text-slate-500">{`Visits: ${short.visits}`}</p>
      </CardFooter>
    </Card>
  );
}
