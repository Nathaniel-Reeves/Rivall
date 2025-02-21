import { Image } from '@/components/ui/image';
import { Box } from '@/components/ui/box';
import { Card } from '@/components/ui/card';
import { VStack } from '@/components/ui/vstack';
import { HStack } from '@/components/ui/hstack';
import { Button, ButtonIcon } from '@/components/ui/button';
import { Text } from '@/components/ui/text';
import { ChevronLeftIcon, AlertCircleIcon } from '@/components/ui/icon';
import { Input, InputField, InputSlot, InputIcon } from "@/components/ui/input"
import {
  FormControl,
  FormControlError,
  FormControlErrorText,
  FormControlErrorIcon,
  FormControlLabel,
  FormControlLabelText,
} from "@/components/ui/form-control"
import { Divider } from "@/components/ui/divider"
import { EyeIcon, EyeOffIcon } from "@/components/ui/icon"
import React from "react"
import { ScrollView } from 'react-native';
import { Link } from 'expo-router';

export default function RegistrationScreen() {

  const [isInvalid, setIsInvalid] = React.useState(false)
  const [email, setEmail] = React.useState("12345")
  const [password, setPassword] = React.useState("12345")
  const [confirmPassword, setConfirmPassword] = React.useState("12345")
  const handleSubmit = () => {
    // Validate Register logic
    if (email.length < 6) {
      setIsInvalid(true)
    } else {
      setIsInvalid(false)
    }
  }
  const [showPassword, setShowPassword] = React.useState(false)
  const handleState = () => {
    setShowPassword((showState) => {
      return !showState
    })
  }

  return (
    <ScrollView
      style={{
        position: 'absolute',
        left: 0,
        right: 0,
        top: 0,
        height: '100%',
      }}
    >
      <Card
        className="align-top m-4 my-10 p-10 shadow-md shadow-black flex-col justify-between"
      >
        <Box>
          <Link href="/welcome" asChild replace>
            <Button
              className="mb-6 w-20"
              variant="outline"
            >
              <ButtonIcon as={ChevronLeftIcon}/>
            </Button>
          </Link>
          <Text
            className="text-typography-800 text-4xl font-normal text-pretty text-left mb-4"
          >
            Hello New Rivall!
          </Text>
          <Text
            className="text-typography-800 text-2xl text-pretty text-left mb-10"
          >
            Let the Challenges Begin!
          </Text>
        </Box>
        <VStack
          className="gap-4 mt-10"
        >
          <FormControl
            isInvalid={isInvalid}
            size="md"
            isDisabled={false}
            isReadOnly={false}
            isRequired={true}
            className="gap-4 mb-4"
          >
            <Box>
              <FormControlLabel>
                <FormControlLabelText>Email</FormControlLabelText>
              </FormControlLabel>
              <Input className="my-1" size="md">
                <InputField
                  type="text"
                  placeholder="password"
                  value={email}
                  onChangeText={(text) => setEmail(text)}
                />
              </Input>
              <FormControlError>
                <FormControlErrorIcon as={AlertCircleIcon} />
                <FormControlErrorText>
                  Atleast 6 characters are required.
                </FormControlErrorText>
              </FormControlError>
            </Box>
            <Box>
              <FormControlLabel>
                <FormControlLabelText>Password</FormControlLabelText>
              </FormControlLabel>
              <Input className="my-1" size="md">
                <InputField
                  type={showPassword ? "text" : "password"}
                  placeholder="password"
                  value={password}
                  onChangeText={(text) => setPassword(text)}
                />
                <InputSlot className="pr-3" onPress={handleState}>
                  <InputIcon as={showPassword ? EyeIcon : EyeOffIcon} />
                </InputSlot>
              </Input>
              <FormControlError>
                <FormControlErrorIcon as={AlertCircleIcon} />
                <FormControlErrorText>
                  Atleast 6 characters are required.
                </FormControlErrorText>
              </FormControlError>
            </Box>
            <Box>
              <FormControlLabel>
                <FormControlLabelText>Confirm Password</FormControlLabelText>
              </FormControlLabel>
              <Input className="my-1" size="md">
                <InputField
                  type={showPassword ? "text" : "password"}
                  placeholder="password"
                  value={confirmPassword}
                  onChangeText={(text) => setConfirmPassword(text)}
                />
                <InputSlot className="pr-3" onPress={handleState}>
                  <InputIcon as={showPassword ? EyeIcon : EyeOffIcon} />
                </InputSlot>
              </Input>
              <FormControlError>
                <FormControlErrorIcon as={AlertCircleIcon} />
                <FormControlErrorText>
                  Atleast 6 characters are required.
                </FormControlErrorText>
              </FormControlError>
            </Box>
          </FormControl>
          <Button
            className="shadow-md shadow-black"
            size="lg"
          >
            <Text
              className="text-typography-0 text-lg"
            >Register</Text>
          </Button>
          {/* <Box className="flex-row justify-center">
            <HStack className="justify-self-center items-center">
              <Divider className="w-12 bg-gray-700" />
              <Text className="text-center mx-3">Or Register With</Text>
              <Divider className="w-12 bg-gray-700" />
            </HStack>
          </Box>
          <Box className="flex-row justify-center">
            <HStack className="flex-1 justify-between">
              <Button variant="outline" size="lg" className="p-8"><Image source={require("@/assets/social-icons/facebook.png")} className="h-[32px] w-[32px]"/></Button>
              <Button variant="outline" size="lg" className="p-8"><Image source={require("@/assets/social-icons/google.png")} className="h-[32px] w-[32px]"/></Button>
              <Button variant="outline" size="lg" className="p-8"><Image source={require("@/assets/social-icons/apple.png")} className="h-[32px] w-[32px]"/></Button>
            </HStack>
          </Box> */}
          <Box className="flex-row justify-center mt-4">
            <HStack className="flex-1 justify-between">
              <Text className="text-center text-lg font-medium">Already have an account?</Text>
              <Link href="/login" asChild replace><Text className="text-info-300 font-semibold text-lg">Login Now</Text></Link>
            </HStack>
          </Box>
        </VStack>
      </Card>
    </ScrollView>
  );
}
