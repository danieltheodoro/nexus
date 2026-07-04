import { useEffect, useRef, useState } from "react";
import {
  Animated,
  Easing,
  Modal,
  Pressable,
  Text,
  TextInput,
  View,
} from "react-native";

export default function CreateAccountScreen() {
  const [confirmVisible, setConfirmVisible] = useState(false);

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
    <View style={{ flex: 1, justifyContent: "center", alignItems: "center", backgroundColor: "#111827", padding: 24 }}>
      <Animated.View style={{
        width: 420,
        maxWidth: "100%",
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
        opacity,
        transform: [{ translateY }],
      }}>
        <Text style={{ fontSize: 38, fontWeight: "700", color: "#FFFFFF", textAlign: "center", marginBottom: 6 }}>
          Nexus
        </Text>

        <Text style={{ fontSize: 16, color: "#9CA3AF", textAlign: "center", marginBottom: 20 }}>
          create your account.
        </Text>

        <Text style={{ color: "#FFFFFF", fontWeight: "500", marginBottom: 5 }}>First name</Text>
        <TextInput placeholder="Enter your first name" placeholderTextColor="#B8C2CF" style={inputStyle} />

        <Text style={labelStyle}>Last name</Text>
        <TextInput placeholder="Enter your last name" placeholderTextColor="#B8C2CF" style={inputStyle} />

        <Text style={labelStyle}>Email</Text>
        <TextInput
          placeholder="Enter your email"
          placeholderTextColor="#B8C2CF"
          keyboardType="email-address"
          autoCapitalize="none"
          style={inputStyle}
        />

        <Text style={labelStyle}>Password</Text>
        <TextInput
          placeholder="Create a password"
          placeholderTextColor="#B8C2CF"
          secureTextEntry
          style={{ ...inputStyle, marginBottom: 20 }}
        />

        <Pressable
          onPress={() => setConfirmVisible(true)}
          style={({ pressed }) => ({
            backgroundColor: pressed ? "#1D4ED8" : "#2563EB",
            borderRadius: 10,
            padding: 16,
            alignItems: "center",
            transform: [{ scale: pressed ? 0.98 : 1 }],
          })}
        >
          <Text style={{ color: "#FFFFFF", fontSize: 16, fontWeight: "600" }}>
            Create account
          </Text>
        </Pressable>

        <View style={{ flexDirection: "row", justifyContent: "center", marginTop: 20 }}>
          <Text style={{ color: "#9CA3AF" }}>Already have an account?</Text>
          <Pressable>
            <Text style={{ color: "#60A5FA", marginLeft: 6, fontWeight: "600" }}>
              Sign In
            </Text>
          </Pressable>
        </View>
      </Animated.View>

      <Modal transparent visible={confirmVisible} animationType="fade">
        <View style={{
          flex: 1,
          backgroundColor: "rgba(0,0,0,0.55)",
          justifyContent: "center",
          alignItems: "center",
          padding: 24,
        }}>
          <View style={{
            width: 360,
            maxWidth: "100%",
            backgroundColor: "#273449",
            borderRadius: 18,
            padding: 28,
            borderWidth: 1,
            borderColor: "#36455E",
          }}>
            <Text style={{ color: "#FFFFFF", fontSize: 24, fontWeight: "700", marginBottom: 8 }}>
              Confirm password
            </Text>

            <Text style={{ color: "#9CA3AF", marginBottom: 20 }}>
              Enter your password again to finish creating your account.
            </Text>

            <Text style={labelStyle}>Password</Text>

            <TextInput
              placeholder="Repeat your password"
              placeholderTextColor="#B8C2CF"
              secureTextEntry
              style={{ ...inputStyle, marginBottom: 20 }}
            />

            <View style={{ flexDirection: "row", gap: 12 }}>
              <Pressable
                onPress={() => setConfirmVisible(false)}
                style={{
                  flex: 1,
                  backgroundColor: "#2B3A51",
                  borderWidth: 1,
                  borderColor: "#4A5E7A",
                  borderRadius: 10,
                  padding: 14,
                  alignItems: "center",
                }}
              >
                <Text style={{ color: "#FFFFFF", fontWeight: "600" }}>Cancel</Text>
              </Pressable>

              <Pressable
                style={{
                  flex: 1,
                  backgroundColor: "#2563EB",
                  borderRadius: 10,
                  padding: 14,
                  alignItems: "center",
                }}
              >
                <Text style={{ color: "#FFFFFF", fontWeight: "600" }}>Confirm</Text>
              </Pressable>
            </View>
          </View>
        </View>
      </Modal>
    </View>
  );
}

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
  marginBottom: 12,
};
