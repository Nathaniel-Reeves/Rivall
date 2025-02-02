import { Button, ButtonText } from "@/components/ui/button"
import { FormControl } from "@/components/ui/form-control"
import { Heading } from "@/components/ui/heading"
import { Input, InputField, InputIcon, InputSlot } from "@/components/ui/input"
import { Text } from "@/components/ui/text"
import { VStack } from "@/components/ui/vstack"
import { HStack } from "@/components/ui/hstack"
import { EyeIcon, EyeOffIcon } from "@/components/ui/icon"
import { useState } from "react"

export default function LoginScreen() {
  const [showPassword, setShowPassword] = useState(false)
  const handleState = () => {
    setShowPassword((showState) => {
      return !showState
    })
  }
  return (
    <FormControl className="p-4 border rounded-lg border-outline-300 max-w-[500px] bg-white m-2">
      <VStack space="xl">
        <Heading className="text-typography-900 leading-3 pt-3">Login</Heading>
        <VStack space="xs">
          <Text className="text-typography-500">Email</Text>
          <Input className="min-w-[250px]">
            <InputField type="text" />
          </Input>
        </VStack>
        <VStack space="xs">
        <Text className="text-typography-500">Password</Text>
        <Input className="text-center">
          <InputField type={showPassword ? "text" : "password"} />
          <InputSlot className="pr-3" onPress={handleState}>
              <InputIcon as={showPassword ? EyeIcon : EyeOffIcon} />
          </InputSlot>
        </Input>
        </VStack>
        <HStack space="sm">
          <Button className="flex-1" variant="outline" onPress={() => {}}
          >
            <ButtonText>Sign Up</ButtonText>
          </Button>
          <Button className="flex-1" onPress={() => {}}
          >
            <ButtonText>Sign In</ButtonText>
          </Button>
        </HStack>
      </VStack>
    </FormControl>
  )
}