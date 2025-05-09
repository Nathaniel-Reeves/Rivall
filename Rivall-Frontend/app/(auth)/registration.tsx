import { Image } from '@/components/ui/image';
import { Box } from '@/components/ui/box';
import { Card } from '@/components/ui/card';
import { VStack } from '@/components/ui/vstack';
import { HStack } from '@/components/ui/hstack';
import { Button, ButtonIcon, ButtonText, ButtonSpinner } from '@/components/ui/button';
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
import { useState, useEffect } from "react"
import { ScrollView } from 'react-native';
import { Link, useRouter } from 'expo-router';

import { register, login } from '@/api/auth';
import { useUserStore } from '@/global-store/user_store';
import {
  validateEmail,
  validateName,
  validatePassword,
  matchPassword,
  validPassword,
  validEmail,
  validName
} from '@/common/auth_helper_functions';

export default function RegistrationScreen() {

  const [isInvalidEmail, setIsInvalidEmail] = useState(false)
  const [accountExists, setAccountExists] = useState(false)
  const [passwordsMatch, setPasswordsMatch] = useState(true)
  const [isInvalidPassword, setIsInvalidPassword] = useState(false)
  const [isInvalidNames, setIsInvalidNames] = useState(false)

  const [firstName, setFirstName] = useState("")
  const [lastName, setLastName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")

  const [registerLoading, setRegisterLoading] = useState(false)

  const setStoreData = useUserStore((state: any) => state.setStoreData)
  const router = useRouter()

  const handleSubmit = async () => {
    setRegisterLoading(true)

    // Validate Login logic
    let flag = false
    if (!validateName(firstName)) {
      setIsInvalidNames(true)
      flag = true
    }

    if (!validateName(lastName)) {
      setIsInvalidNames(true)
      flag = true
    }

    if (!validateEmail(email)) {
      setIsInvalidEmail(true)
      flag = true
    }

    if (!validatePassword(password)) {
      setIsInvalidPassword(true)
      flag = true
    }

    if (!validatePassword(confirmPassword)) {
      setIsInvalidPassword(true)
      flag = true
    }

    if (!matchPassword(password, confirmPassword)) {
      setPasswordsMatch(true)
      flag = true
    }

    if (flag) {
      setRegisterLoading(false)
      return
    }

    // Send Register Request
    const res1 = await register(firstName, lastName, email, password)
    if (res1 === null) {
      console.debug('Registration Failed')
      setAccountExists(res1.status === 403)
      setRegisterLoading(false)
      return
    }

    const res2 = await login(email, password)
    if (res2 === null) {
      console.debug('Login Failed')
      setRegisterLoading(false)
      return
    }

    console.debug('Login Successful')
    setStoreData(res2.data)
    setRegisterLoading(false)
    router.replace('/entry/home')
    return
  }

  const [showPassword, setShowPassword] = useState(false)
  const handleState = () => {
    setShowPassword((showState) => {
      return !showState
    })
  }

  useEffect(() => {
    setIsInvalidNames(!validName(firstName) || !validName(lastName))
  }, [firstName, lastName])

  useEffect(() => {
    setPasswordsMatch(matchPassword(password, confirmPassword))
    setIsInvalidPassword(!validPassword(password) || !validPassword(confirmPassword))
  }, [password, confirmPassword])

  useEffect(() => {
    setIsInvalidEmail(!validEmail(email))
  }, [email])

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
        className="align-top m-4 my-16 p-10 shadow-md shadow-black flex-col justify-between"
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
            isInvalid={isInvalidNames}
            size="md"
            isDisabled={false}
            isReadOnly={false}
            isRequired={true}
            className="gap-4"
          >
            <Box>
              <FormControlLabel>
                <FormControlLabelText>First Name</FormControlLabelText>
              </FormControlLabel>
              <Input className="my-1" size="md">
                <InputField
                  type="text"
                  placeholder="first name"
                  value={firstName}
                  onChangeText={(text) => setFirstName(text)}
                />
              </Input>
            </Box>
            <Box>
              <FormControlLabel>
                <FormControlLabelText>Last Name</FormControlLabelText>
              </FormControlLabel>
              <Input className="my-1" size="md">
                <InputField
                  type="text"
                  placeholder="first name"
                  value={lastName}
                  onChangeText={(text) => setLastName(text)}
                />
              </Input>
              <FormControlError>
                <FormControlErrorIcon as={AlertCircleIcon} />
                <FormControlErrorText>
                First and last name are required.
                </FormControlErrorText>
              </FormControlError>
            </Box>
          </FormControl>
          <FormControl
            isInvalid={isInvalidEmail || accountExists}
            size="md"
            isDisabled={false}
            isReadOnly={false}
            isRequired={true}
            className="gap-4"
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
                {isInvalidEmail ? "Invalid Email.  " : ""}
                {accountExists ? "Account already exists." : ""}
                </FormControlErrorText>
              </FormControlError>
            </Box>
          </FormControl>
          <FormControl
            isInvalid={!passwordsMatch || isInvalidPassword}
            size="md"
            isDisabled={false}
            isReadOnly={false}
            isRequired={true}
            className="gap-4 mb-4"
          >
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
                  {isInvalidPassword ? "Invalid Password.  " : ""}
                  {passwordsMatch ? "" : "Passwords do not match."}
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
                  placeholder="confirm password"
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
                  {passwordsMatch ? "" : "Passwords do not match."}
                  {isInvalidPassword ? "Invalid Password.  " : ""}
                </FormControlErrorText>
              </FormControlError>
            </Box>
          </FormControl>
          <Button
            className="shadow-md shadow-black"
            size="lg"
            onPress={handleSubmit}
            disabled={registerLoading}
          >
            {registerLoading ? <ButtonSpinner/> : <ButtonText className="text-typography-0 text-lg">Register</ButtonText>}
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
