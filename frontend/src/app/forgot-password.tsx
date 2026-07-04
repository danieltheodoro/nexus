import { useEffect, useRef, useState } from "react";
import {
  Animated,
  Easing,
  Pressable,
  Text,
  TextInput,
  View,
} from "react-native";

export default function ForgotPasswordScreen() {
  const [sent, setSent] = useState(false);

  const translateY = useRef(new Animated.Value(20)).current;
  const opacity = useRef(new Animated.Value(0)).current;

  useEffect(() => {
    Animated.parallel([
      Animated.spring(translateY, {
        toValue: 0,
        damping: 12,
        stiffness: 120,
        mass: 0.8,
        useNativeDriver: true,
      }),
      Animated.timing(opacity, {
        toValue: 1,
        duration: 450,
        easing: Easing.out(Easing.cubic),
        useNativeDriver: true,
      }),
    ]).start();
  }, []);

  return (
    <View style={screenStyle}>
      <Animated.View style={{ ...cardStyle, opacity, transform: [{ translateY }] }}>
        <Text style={titleStyle}>Nexus</Text>

        <Text style={subtitleStyle}>
          {sent ? "check your inbox." : "reset your password."}
        </Text>

        {!sent ? (
          <>
            <Text style={labelStyle}>Email</Text>

            <TextInput
              placeholder="Enter your email"
              placeholderTextColor="#B8C2CF"
              keyboardType="email-address"
              autoCapitalize="none"
              style={{ ...inputStyle, marginBottom: 20 }}
            />

            <Pressable onPress={() => setSent(true)} style={primaryButtonStyle}>
              <Text style={buttonTextStyle}>Send reset link</Text>
            </Pressable>

            <View style={dividerStyle}>
              <View style={lineStyle} />
              <Text style={{ color: "#9CA3AF", marginHorizontal: 12 }}>or</Text>
              <View style={lineStyle} />
            </View>

            <Pressable style={secondaryButtonStyle}>
              <Text style={buttonTextStyle}>Back to Sign In</Text>
            </Pressable>
          </>
        ) : (
          <>
            <Text
              style={{
                color: "#B8C2CF",
                textAlign: "center",
                lineHeight: 22,
                marginBottom: 24,
              }}
            >
              We sent a password reset link to your email.
            </Text>

            <Pressable style={primaryButtonStyle}>
              <Text style={buttonTextStyle}>Send again</Text>
            </Pressable>

            <Pressable onPress={() => setSent(false)}>
              <Text
                style={{
                  color: "#60A5FA",
                  textAlign: "center",
                  marginTop: 20,
                  fontWeight: "600",
                }}
              >
                Use another email
              </Text>
            </Pressable>
          </>
        )}
      </Animated.View>
    </View>
  );
}

const screenStyle = {
  flex: 1,
  justifyContent: "center" as const,
  alignItems: "center" as const,
  backgroundColor: "#111827",
  padding: 24,
};

const cardStyle = {
  width: 420,
  maxWidth: "100%" as const,
  backgroundColor: "#273449",
  borderRadius: 20,
  paddingHorizontal: 32,
  paddingTop: 24,
  paddingBottom: 28,
  borderWidth: 1,
  borderColor: "#36455E",
  shadowColor: "#000",
  shadowOpacity: 0.35,
  shadowRadius: 24,
  shadowOffset: { width: 0, height: 10 },
  elevation: 12,
};

const titleStyle = {
  fontSize: 38,
  fontWeight: "700" as const,
  color: "#FFFFFF",
  textAlign: "center" as const,
  marginBottom: 6,
};

const subtitleStyle = {
  fontSize: 16,
  color: "#9CA3AF",
  textAlign: "center" as const,
  marginBottom: 20,
};

const labelStyle = {
  color: "#FFFFFF",
  fontWeight: "500" as const,
  marginBottom: 5,
};

const inputStyle = {
  backgroundColor: "#42546F",
  color: "#FFFFFF",
  borderWidth: 1,
  borderColor: "#5B708D",
  borderRadius: 10,
  paddingHorizontal: 16,
  paddingVertical: 16,
};

const primaryButtonStyle = ({ pressed }: { pressed: boolean }) => ({
  backgroundColor: pressed ? "#1D4ED8" : "#2563EB",
  borderRadius: 10,
  padding: 16,
  alignItems: "center" as const,
  transform: [{ scale: pressed ? 0.98 : 1 }],
});

const secondaryButtonStyle = ({ pressed }: { pressed: boolean }) => ({
  backgroundColor: pressed ? "#334764" : "#2B3A51",
  borderWidth: 1,
  borderColor: "#4A5E7A",
  borderRadius: 10,
  padding: 16,
  alignItems: "center" as const,
  transform: [{ scale: pressed ? 0.98 : 1 }],
});

const buttonTextStyle = {
  color: "#FFFFFF",
  fontSize: 16,
  fontWeight: "600" as const,
};

const dividerStyle = {
  flexDirection: "row" as const,
  alignItems: "center" as const,
  marginVertical: 18,
};

const lineStyle = {
  flex: 1,
  height: 1,
  backgroundColor: "#4A5E7A",
};
