import { Button, ButtonText } from "@/components/ui/button"
import { FormControl } from "@/components/ui/form-control"
import { Heading } from "@/components/ui/heading"
import { Input, InputField, InputIcon, InputSlot } from "@/components/ui/input"
import { Text } from "@/components/ui/text"
import { VStack } from "@/components/ui/vstack"
import { HStack } from "@/components/ui/hstack"
import { EyeIcon, EyeOffIcon } from "@/components/ui/icon"
import { useState } from "react"
import { Redirect } from "expo-router"
import { useMutation } from "@tanstack/react-query"

import { login, register } from '@/api/auth';
import { useAuth } from "@/global-store/authStore"

export default function LoginScreen() {
  const [showPassword, setShowPassword] = useState(false)

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  const setUser = useAuth((state) => state.setUser);
  const setToken = useAuth((state) => state.setToken);
  const isLoggedIn = useAuth((state) => !!state.token);

  const loginMutation = useMutation({ 
    mutationFn: () => login(email, password),
    onSuccess: (data) => {
      console.log('login success', data);
      if (data.user && data.token) {
        setUser(data.user);
        setToken(data.token);
      }
    }, 
    onError: (error) => {
      console.log('login error:', error);
    }
  });

  const registerMutation = useMutation({
    mutationFn: () => register(email, password),
    onSuccess: (data) => {
      console.log('register success', data);
      if (data.user && data.token) {
        setUser(data.user);
        setToken(data.token);
      }
    },
    onError: (error) => {
      console.log('register error:', error);
    }
  });

  const handleState = () => {
    setShowPassword((showState) => {
      return !showState
    })
  }

  if (isLoggedIn) {
    return <Redirect href="/" />
  }

  return (
    <FormControl isInvalid={loginMutation.error || registerMutation.error} className="p-4 border rounded-lg border-outline-300 max-w-[500px] bg-white m-2">
      <VStack space="xl">
        <Heading className="text-typography-900 leading-3 pt-3">Login</Heading>
        <VStack space="xs">
          <Text className="text-typography-500">Email</Text>
          <Input className="min-w-[250px]">
            <InputField type="text" value={email} onChangeText={setEmail} />
          </Input>
        </VStack>
        <VStack space="xs">
        <Text className="text-typography-500">Password</Text>
        <Input className="text-center">
          <InputField type={showPassword ? "text" : "password"} value={password} onChangeText={setPassword}/>
          <InputSlot className="pr-3" onPress={handleState}>
              <InputIcon as={showPassword ? EyeIcon : EyeOffIcon} />
          </InputSlot>
        </Input>
        </VStack>
        <HStack space="sm">
          <Button className="flex-1" variant="outline" onPress={() => registerMutation.mutate() }
          >
            <ButtonText>Sign Up</ButtonText>
          </Button>
          <Button className="flex-1" onPress={() => loginMutation.mutate() }
          >
            <ButtonText>Sign In</ButtonText>
          </Button>
        </HStack>
      </VStack>
    </FormControl>
  )
}