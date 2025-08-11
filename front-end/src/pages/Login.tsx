import { SERVER_URL } from "@/utils/constants";
import { Button } from "@heroui/button";
import { Card, CardHeader, CardBody } from "@heroui/card";
import { Divider } from "@heroui/divider";
import { GoogleLogo } from "phosphor-react";

export function Login() {
  const handleLogin = () => {
    window.location.href = `${SERVER_URL}/auth/google/login`;
  };

  return (
    <div className="w-full h-screen flex justify-center items-center">
      <Card
        shadow="md"
        className="min-w-[400px] max-h-[320px] flex flex-col items-center px-12 py-8 space-y-8 "
      >
        <CardHeader className="text-center flex flex-col items-center space-y-4">
          <h1 className="text-4xl font-semibold text-green-500">MinUrly</h1>
          <Divider orientation="horizontal" />
        </CardHeader>
        <CardBody className="flex justify-center items-center">
          <Button
            color="primary"
            radius="sm"
            startContent={<GoogleLogo size={25} />}
            onPress={handleLogin}
            className="w-full font-semibold text-lg"
          >
            Login with Google
          </Button>
        </CardBody>
      </Card>
    </div>
  );
}
