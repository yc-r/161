package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	_ "github.com/google/uuid"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	// var doris *client.User
	// var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	// dorisFile := "dorisFile.txt"
	// eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

        // 新增测试：重复 InitUser 应返回错误
        Specify("Extra Test: Duplicate InitUser should return an error.", func() {
            userlib.DebugMsg("Initializing user dupuser.")
            u1, err := client.InitUser("dupuser", "pw123")
            Expect(err).To(BeNil())
            Expect(u1).ToNot(BeNil())

            // 第二次用相同用户名初始化应失败
            _, err = client.InitUser("dupuser", "pw123")
            Expect(err).ToNot(BeNil())

            // 用正确凭据 GetUser 应该仍然可以
            u2, err := client.GetUser("dupuser", "pw123")
            Expect(err).To(BeNil())
            Expect(u2).ToNot(BeNil())
        })

        // 新增测试：撤销应仅影响被撤销的文件（隔离性）
        Specify("Extra Test: Revoke should only affect the revoked file (isolation test).", func() {
            userlib.DebugMsg("Initializing users AliceX and BobX.")
            aliceX, err := client.InitUser("aliceX", "pwA")
            Expect(err).To(BeNil())
            bobX, err := client.InitUser("bobX", "pwB")
            Expect(err).To(BeNil())

            // AliceX 存两个文件
            err = aliceX.StoreFile("fileA", []byte("AAA"))
            Expect(err).To(BeNil())
            err = aliceX.StoreFile("fileB", []byte("BBB"))
            Expect(err).To(BeNil())

            // 分享给 BobX
            inviteA, err := aliceX.CreateInvitation("fileA", "bobX")
            Expect(err).To(BeNil())
            err = bobX.AcceptInvitation("aliceX", inviteA, "bobFileA")
            Expect(err).To(BeNil())

            inviteB, err := aliceX.CreateInvitation("fileB", "bobX")
            Expect(err).To(BeNil())
            err = bobX.AcceptInvitation("aliceX", inviteB, "bobFileB")
            Expect(err).To(BeNil())

            // BobX 能读取两个文件
            data, err := bobX.LoadFile("bobFileA")
            Expect(err).To(BeNil())
            Expect(string(data)).To(Equal("AAA"))
            data, err = bobX.LoadFile("bobFileB")
            Expect(err).To(BeNil())
            Expect(string(data)).To(Equal("BBB"))

            // AliceX 撤销 BobX 对 fileA 的访问
            err = aliceX.RevokeAccess("fileA", "bobX")
            Expect(err).To(BeNil())

            // BobX 不应能访问 bobFileA，但应仍能访问 bobFileB
            _, err = bobX.LoadFile("bobFileA")
            Expect(err).ToNot(BeNil())

            data, err = bobX.LoadFile("bobFileB")
            Expect(err).To(BeNil())
            Expect(string(data)).To(Equal("BBB"))
        })

        // 新增测试组：更多边界与异常情况
        Specify("Extra Test: Corrupted invite should not be accepted.", func() {
            alice, err := client.InitUser("aliceCorrupt", "pw1")
            Expect(err).To(BeNil())
            bob, err := client.InitUser("bobCorrupt", "pw2")
            Expect(err).To(BeNil())

            err = alice.StoreFile("secret", []byte("TOPSECRET"))
            Expect(err).To(BeNil())

            invite, err := alice.CreateInvitation("secret", "bobCorrupt")
            Expect(err).To(BeNil())
            // 模拟损坏 invite（复制并修改字节）
            corrupt := make([]byte, len(invite))
            copy(corrupt, invite)
            if len(corrupt) > 0 {
                corrupt[0] ^= 0xFF
            }
            // Bob 尝试接受被损坏的 invite 应失败
            err = bob.AcceptInvitation("aliceCorrupt", corrupt, "bobSecret")
            Expect(err).ToNot(BeNil())
        })

        Specify("Extra Test: Multiple invites to same recipient can be accepted under different names.", func() {
            alice, err := client.InitUser("aliceMulti", "pwA")
            Expect(err).To(BeNil())
            bob, err := client.InitUser("bobMulti", "pwB")
            Expect(err).To(BeNil())

            err = alice.StoreFile("sharedFile", []byte("DATA"))
            Expect(err).To(BeNil())

            inv1, err := alice.CreateInvitation("sharedFile", "bobMulti")
            Expect(err).To(BeNil())
            inv2, err := alice.CreateInvitation("sharedFile", "bobMulti")
            Expect(err).To(BeNil())

            // Bob 接受两次 invite，保存为不同名字
            err = bob.AcceptInvitation("aliceMulti", inv1, "bobCopy1")
            Expect(err).To(BeNil())
            err = bob.AcceptInvitation("aliceMulti", inv2, "bobCopy2")
            Expect(err).To(BeNil())

            // Bob 可以从两个名字读取同样内容
            d1, err := bob.LoadFile("bobCopy1")
            Expect(err).To(BeNil())
            d2, err := bob.LoadFile("bobCopy2")
            Expect(err).To(BeNil())
            Expect(string(d1)).To(Equal("DATA"))
            Expect(string(d2)).To(Equal("DATA"))
        })

        Specify("Extra Test: Revoke prevents future appends by revoked user.", func() {
            alice, err := client.InitUser("aliceRevoke", "pwa")
            Expect(err).To(BeNil())
            bob, err := client.InitUser("bobRevoke", "pwb")
            Expect(err).To(BeNil())

            err = alice.StoreFile("log", []byte("start"))
            Expect(err).To(BeNil())

            invite, err := alice.CreateInvitation("log", "bobRevoke")
            Expect(err).To(BeNil())
            err = bob.AcceptInvitation("aliceRevoke", invite, "bobLog")
            Expect(err).To(BeNil())

            // Bob 能 append
            err = bob.AppendToFile("bobLog", []byte("-ok"))
            Expect(err).To(BeNil())

            // Alice revoke Bob
            err = alice.RevokeAccess("log", "bobRevoke")
            Expect(err).To(BeNil())

            // Bob 尝试 append 应失败
            err = bob.AppendToFile("bobLog", []byte("-later"))
            Expect(err).ToNot(BeNil())

            // Alice 仍能读取 original file
            data, err := alice.LoadFile("log")
            Expect(err).To(BeNil())
            Expect(string(data)).To(ContainSubstring("start"))
        })
	})
})
