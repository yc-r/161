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
        Specify("Confidentiality Test: uninvited user cannot read file", func() {
            alice, err := client.InitUser("aliceConf", "pwA")
            Expect(err).To(BeNil())
            bob, err := client.InitUser("bobConf", "pwB")
            Expect(err).To(BeNil())
            charlie, err := client.InitUser("charlieConf", "pwC")
            Expect(err).To(BeNil())

            // Alice 存文件并只分享给 Bob
            err = alice.StoreFile("topsecret", []byte("SENSITIVE"))
            Expect(err).To(BeNil())

            invite, err := alice.CreateInvitation("topsecret", "bobConf")
            Expect(err).To(BeNil())
            err = bob.AcceptInvitation("aliceConf", invite, "bobSecret")
            Expect(err).To(BeNil())

            // Charlie 没有被邀请，不能读取
            _, err = charlie.LoadFile("bobSecret")
            Expect(err).ToNot(BeNil())
            _, err = charlie.LoadFile("topsecret")
            Expect(err).ToNot(BeNil())
        })

        // 额外测试：机密性 - 针对某个接收者的 invite 不应被其它用户接受
        Specify("Confidentiality Test: invite for Bob cannot be accepted by Charlie", func() {
            alice, err := client.InitUser("aliceInvite", "pw1")
            Expect(err).To(BeNil())
            bob, err := client.InitUser("bobInvite", "pw2")
            Expect(err).To(BeNil())
            charlie, err := client.InitUser("charlieInvite", "pw3")
            Expect(err).To(BeNil())

            err = alice.StoreFile("doc", []byte("DOC"))
            Expect(err).To(BeNil())

            invite, err := alice.CreateInvitation("doc", "bobInvite")
            Expect(err).To(BeNil())

            // Charlie 尝试接受发给 Bob 的 invite 应失败
            err = charlie.AcceptInvitation("aliceInvite", invite, "charlieDoc")
            Expect(err).ToNot(BeNil())

            // Bob 能正常接受
            err = bob.AcceptInvitation("aliceInvite", invite, "bobDoc")
            Expect(err).To(BeNil())
        })

        // 额外测试：完整性 - 重放/再次使用已被撤销或已使用的 invite 应被拒绝（若实现了防护）
        Specify("Integrity Test: replaying or reusing an invite after revoke/accept should not grant extra access", func() {
            alice, err := client.InitUser("aliceReplay", "pwA")
            Expect(err).To(BeNil())
            bob, err := client.InitUser("bobReplay", "pwB")
            Expect(err).To(BeNil())

            err = alice.StoreFile("journal", []byte("LOG"))
            Expect(err).To(BeNil())

            invite, err := alice.CreateInvitation("journal", "bobReplay")
            Expect(err).To(BeNil())

            // Bob 接受 invite（第一次接受应该成功）
            err = bob.AcceptInvitation("aliceReplay", invite, "bobJournal1")
            Expect(err).To(BeNil())

            // Alice 撤销 Bob 的访问
            err = alice.RevokeAccess("journal", "bobReplay")
            Expect(err).To(BeNil())

            // Bob 再次尝试用同一个 invite 接受为另一个名字，应被拒绝（如果实现了防护）
            err = bob.AcceptInvitation("aliceReplay", invite, "bobJournal2")
            Expect(err).ToNot(BeNil())
        })

        // 额外测试：完整性 - 被接受的邀请所产生的修改对所有被授权方可见
        Specify("Integrity Test: updates from authorized users are visible to owner and other authorized users", func() {
            alice, err := client.InitUser("aliceInt", "pwa")
            Expect(err).To(BeNil())
            bob, err := client.InitUser("bobInt", "pwb")
            Expect(err).To(BeNil())

            err = alice.StoreFile("shared", []byte("INIT"))
            Expect(err).To(BeNil())

            invite, err := alice.CreateInvitation("shared", "bobInt")
            Expect(err).To(BeNil())
            err = bob.AcceptInvitation("aliceInt", invite, "bobShared")
            Expect(err).To(BeNil())

            // Bob append，Alice 应能看到改动
            err = bob.AppendToFile("bobShared", []byte("-BOB"))
            Expect(err).To(BeNil())

            data, err := alice.LoadFile("shared")
            Expect(err).To(BeNil())
            Expect(string(data)).To(ContainSubstring("-BOB"))
        })
        Specify("Extra Test: owner can re-share after revoke and recipient can regain access", func() {
            alice, err := client.InitUser("aliceReshare", "pwA")
            Expect(err).To(BeNil())
            bob, err := client.InitUser("bobReshare", "pwB")
            Expect(err).To(BeNil())

            // Alice 存文件并分享给 Bob
            err = alice.StoreFile("doc", []byte("FIRST"))
            Expect(err).To(BeNil())

            invite, err := alice.CreateInvitation("doc", "bobReshare")
            Expect(err).To(BeNil())
            err = bob.AcceptInvitation("aliceReshare", invite, "bobDoc1")
            Expect(err).To(BeNil())

            // Bob 能读
            data, err := bob.LoadFile("bobDoc1")
            Expect(err).To(BeNil())
            Expect(string(data)).To(Equal("FIRST"))

            // Alice revoke Bob
            err = alice.RevokeAccess("doc", "bobReshare")
            Expect(err).To(BeNil())

            // Bob 不能读
            _, err = bob.LoadFile("bobDoc1")
            Expect(err).ToNot(BeNil())

            // Alice 重新创建 invite 并分享，Bob 接受后应能再次访问
            invite2, err := alice.CreateInvitation("doc", "bobReshare")
            Expect(err).To(BeNil())
            err = bob.AcceptInvitation("aliceReshare", invite2, "bobDoc2")
            Expect(err).To(BeNil())

            data, err = bob.LoadFile("bobDoc2")
            Expect(err).To(BeNil())
            Expect(string(data)).To(Equal("FIRST"))
        })

        // 额外测试：错误的 owner 名称接受 invite 应失败（防止被替代的邀请）
        Specify("Extra Test: AcceptInvitation with wrong owner name should fail", func() {
            alice, err := client.InitUser("aliceWrongOwner", "pwA")
            Expect(err).To(BeNil())
            bob, err := client.InitUser("bobWrongOwner", "pwB")
            Expect(err).To(BeNil())

            err = alice.StoreFile("fileX", []byte("XDATA"))
            Expect(err).To(BeNil())

            invite, err := alice.CreateInvitation("fileX", "bobWrongOwner")
            Expect(err).To(BeNil())

            // Bob 尝试用错误的 owner 名称接受，应失败
            err = bob.AcceptInvitation("aliceWrong", invite, "bobX")
            Expect(err).ToNot(BeNil())

            // 正确 owner 名称接受应成功
            err = bob.AcceptInvitation("aliceWrongOwner", invite, "bobX2")
            Expect(err).To(BeNil())
            data, err := bob.LoadFile("bobX2")
            Expect(err).To(BeNil())
            Expect(string(data)).To(Equal("XDATA"))
        })

        // 额外测试：不同用户同名文件互不干扰（文件名隔离）
        Specify("Extra Test: filename isolation between users", func() {
            alice, err := client.InitUser("aliceIso", "pwA")
            Expect(err).To(BeNil())
            bob, err := client.InitUser("bobIso", "pwB")
            Expect(err).To(BeNil())

            err = alice.StoreFile("sharedName", []byte("ALICE"))
            Expect(err).To(BeNil())
            err = bob.StoreFile("sharedName", []byte("BOB"))
            Expect(err).To(BeNil())

            da, err := alice.LoadFile("sharedName")
            Expect(err).To(BeNil())
            Expect(string(da)).To(Equal("ALICE"))

            db, err := bob.LoadFile("sharedName")
            Expect(err).To(BeNil())
            Expect(string(db)).To(Equal("BOB"))
        })

        // 额外测试：并发会话交错 append 的可见性（简单顺序一致性）
        Specify("Extra Test: interleaved appends from multiple sessions are visible to all", func() {
            aliceMain, err := client.InitUser("aliceInter", "pwA")
            Expect(err).To(BeNil())
            // 第二个会话
            aliceOther, err := client.GetUser("aliceInter", "pwA")
            Expect(err).To(BeNil())

            err = aliceMain.StoreFile("log", []byte("start"))
            Expect(err).To(BeNil())

            // 两个会话交替 append
            err = aliceMain.AppendToFile("log", []byte("-A1"))
            Expect(err).To(BeNil())
            err = aliceOther.AppendToFile("log", []byte("-A2"))
            Expect(err).To(BeNil())
            err = aliceMain.AppendToFile("log", []byte("-A3"))
            Expect(err).To(BeNil())

            data, err := aliceOther.LoadFile("log")
            Expect(err).To(BeNil())
            Expect(string(data)).To(ContainSubstring("start"))
            Expect(string(data)).To(ContainSubstring("-A1"))
            Expect(string(data)).To(ContainSubstring("-A2"))
            Expect(string(data)).To(ContainSubstring("-A3"))
        })

        // 额外测试：错误密码 GetUser 应该失败
        Specify("Auth Test: GetUser with wrong password should fail", func() {
            _, err := client.InitUser("authUser", "rightpw")
            Expect(err).To(BeNil())

            // 用错误密码 GetUser 应返回错误
            _, err = client.GetUser("authUser", "wrongpw")
            Expect(err).ToNot(BeNil())

            // 用正确密码应成功
            u, err := client.GetUser("authUser", "rightpw")
            Expect(err).To(BeNil())
            Expect(u).ToNot(BeNil())
        })

        // 额外测试：StoreFile 覆盖行为（再次 Store 覆盖原内容）
        Specify("Storage Test: StoreFile overwrites existing file", func() {
            alice, err := client.InitUser("aliceOverwrite", "pw")
            Expect(err).To(BeNil())

            err = alice.StoreFile("notes", []byte("first"))
            Expect(err).To(BeNil())

            // 覆盖
            err = alice.StoreFile("notes", []byte("second"))
            Expect(err).To(BeNil())

            data, err := alice.LoadFile("notes")
            Expect(err).To(BeNil())
            Expect(string(data)).To(Equal("second"))
        })

        // 额外测试：同一 invite 被同一接受者、同一名字接收第二次应失败
        Specify("Invite Test: accepting same invite twice under same name should fail", func() {
            alice, err := client.InitUser("aliceReuse", "pwA")
            Expect(err).To(BeNil())
            bob, err := client.InitUser("bobReuse", "pwB")
            Expect(err).To(BeNil())

            err = alice.StoreFile("log", []byte("L"))
            Expect(err).To(BeNil())

            invite, err := alice.CreateInvitation("log", "bobReuse")
            Expect(err).To(BeNil())

            // 第一次接受应成功
            err = bob.AcceptInvitation("aliceReuse", invite, "bobLog")
            Expect(err).To(BeNil())

            // 第二次用同样的 invite 和同样名字接受应失败
            err = bob.AcceptInvitation("aliceReuse", invite, "bobLog")
            Expect(err).ToNot(BeNil())
        })

        // 额外测试：Owner 不应能用 AcceptInvitation 接受自己发出的 invite（应报错）
        Specify("Invite Test: owner cannot accept their own invitation", func() {
            alice, err := client.InitUser("aliceSelf", "pwA")
            Expect(err).To(BeNil())

            err = alice.StoreFile("doc", []byte("D"))
            Expect(err).To(BeNil())

            invite, err := alice.CreateInvitation("doc", "aliceSelf")
            Expect(err).To(BeNil())

            // 尝试让 owner 自己接受应该被拒绝
            err = alice.AcceptInvitation("aliceSelf", invite, "aliceDoc")
            Expect(err).ToNot(BeNil())
        })
	})
})
