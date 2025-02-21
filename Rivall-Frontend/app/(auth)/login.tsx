import { LinearGradient } from 'expo-linear-gradient';
// import { Image } from '@/components/ui/image';
import { Box } from '@/components/ui/box';
import { Card } from '@/components/ui/card';
import { VStack } from '@/components/ui/vstack';
import { HStack } from '@/components/ui/hstack';
import { Button, ButtonIcon, ButtonSpinner, ButtonText } from '@/components/ui/button';
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
  FormControlHelper,
  FormControlHelperText,
} from "@/components/ui/form-control"
// import { Divider } from "@/components/ui/divider"
import { EyeIcon, EyeOffIcon } from "@/components/ui/icon"
import { useState }from "react"
import { ScrollView } from 'react-native';
import { Link, useRouter } from 'expo-router';

import { login } from '@/api/auth';
import { useCredentials } from '@/global-store/credential_store';

function validateEmail(email: string) : boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

function validatePassword(password: string) : boolean {
  return password.length >= 6
}

async function loginUser(email: string, password: string) : Promise<boolean> {
  // Send Login Request
  const res = await login(email, password);
  const userData = res[0];
  const status = res[1];

  if (userData.logged_in && status === 200) {
    // Save User Data
    const setCredentials = useCredentials((state: any) => state.setCredentials);
    setCredentials({ userID: userData.id, token: userData.token })
    return true
  }
  return false
}

export default function LoginScreen() {

  const [isInvalid, setIsInvalid] = useState(false)
  const [loginLoading, setLoginLoading] = useState(false)
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")

  const router = useRouter()

  const handleSubmit = async () => {
    setLoginLoading(true)

    // Validate Login logic
    if (!validateEmail(email)) {
      setIsInvalid(true)
      setLoginLoading(false)
      return
    }

    if (!validatePassword(password)) {
      setIsInvalid(true)
      setLoginLoading(false)
      return
    }

    // Send Login Request
    const result = await loginUser(email, password)
    if (result === false) {
      setIsInvalid(true)
      setLoginLoading(false)
      return
    }

    setLoginLoading(false)
    router.replace('/entry')
  }


  const [showPassword, setShowPassword] = useState(false)
  const handleState = () => {
    setShowPassword((showState) => {
      return !showState
    })
  }

  return (
    <LinearGradient
      // Vertical Background Linear Gradient
      colors={['#77FBFF', '#26C1FE']}
      start={[0, 0]}
      end={[0, 1]}
      style={{
        position: 'absolute',
        left: 0,
        right: 0,
        top: 0,
        height: '100%',
      }}
    >
      <ScrollView
        style={{
          position: 'absolute',
          left: 0,
          right: 0,
          top: 0,
          height: '100%',
        }}
      >
        <Card className="align-top m-4 my-10 p-10 shadow-md shadow-black flex-col justify-between">
          <Box>
            <Link href="/welcome" asChild replace>
              <Button
                className="mb-6 w-20"
                variant="outline"
              >
                <ButtonIcon as={ChevronLeftIcon}/>
              </Button>
            </Link>
            <Text className="text-typography-800 text-4xl font-normal text-pretty text-left mb-4">Welcome Back!</Text>
            <Text className="text-typography-800 text-2xl text-pretty text-left mb-10">Let the Challenges Begin!</Text>
          </Box>
          <VStack className="gap-4 mt-10">
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
                    placeholder="email"
                    value={email}
                    onChangeText={(text) => setEmail(text.toLowerCase())}
                  />
                </Input>
                <FormControlError>
                  <FormControlErrorIcon as={AlertCircleIcon} />
                  <FormControlErrorText>
                    Invalid Email.
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
                    Invalid Password.
                  </FormControlErrorText>
                </FormControlError>
                <FormControlHelper>
                  <FormControlHelperText className="text-right text-typography-500 font-medium">
                    Forgot Password?
                  </FormControlHelperText>
                </FormControlHelper>
              </Box>
            </FormControl>
            <Button
              className="shadow-md shadow-black"
              size="lg"
              onPress={handleSubmit}
              disabled={loginLoading}
            >
              {loginLoading ? <ButtonSpinner/> : <ButtonText className="text-typography-0 text-lg">Login</ButtonText>}
            </Button>
            {/* <Box className="flex-row justify-center">
              <HStack className="justify-self-center items-center">
                <Divider className="w-12 bg-gray-700" />
                <Text className="text-center mx-3">Or Login With</Text>
                <Divider className="w-12 bg-gray-700" />
              </HStack>
            </Box>
            <Box className="flex-row justify-center">
              <HStack className="flex-1 justify-between">
                <Button variant="outline" size="lg" className="p-8" disabled={true}><Image source={require("@/assets/social-icons/facebook.png")} size="xs"/></Button>
                <Button variant="outline" size="lg" className="p-8" disabled={true}><Image source={require("@/assets/social-icons/google.png")} size="xs"/></Button>
                <Button variant="outline" size="lg" className="p-8" disabled={true}><Image source={require("@/assets/social-icons/apple.png")} size="xs"/></Button>
              </HStack>
            </Box> */}
            <Box className="flex-row justify-center mt-4">
              <HStack className="flex-1 justify-between">
                <Text className="text-center text-lg font-medium">Don't have an account?</Text>
                <Link href="/registration" asChild replace><Text className="text-info-300 font-semibold text-lg">Register Now</Text></Link>
              </HStack>
            </Box>
          </VStack>
        </Card>
      </ScrollView>
    </LinearGradient>
  );
}
