package gost341112

import (
	"encoding/hex"
	"fmt"
	"io"
	"testing"
)

type gost3411_2012Test struct {
	out string
	in  string
}

var standard = []gost3411_2012Test{
	{"1b54d01a4af5b9d5cc3d86d68d285462b19abc2475222f35c085122be4ba1ffa00ad30f8767b3a82384c6574f024c311e2a481332b08ef7f41797891c1646f48", "303132333435363738393031323334353637383930313233343536373839303132333435363738393031323334353637383930313233343536373839303132"},
	{"1e88e62226bfca6f9994f1f2d51569e0daf8475a3b0fe61a5300eee46d961376035fe83549ada2b8620fcd7c496ce5b33f0cb9dddc2b6460143b03dabac9fb28", "d1e520e2e5f2f0e82c20d1f2f0e8e1eee6e820e2edf3f6e82c20e2e5fef2fa20f120eceef0ff20f1f2f0e5ebe0ece820ede020f5f0e0e1f0fbff20efebfaeafb20c8e3eef0e5e2fb"},
}

var standard256 = []gost3411_2012Test{
	{"9d151eefd8590b89daa6ba6cb74af9275dd051026bb149a452fd84e5e57b5500", "303132333435363738393031323334353637383930313233343536373839303132333435363738393031323334353637383930313233343536373839303132"},
	{"9dd2fe4e90409e5da87f53976d7405b0c0cac628fc669a741d50063c557e8f50", "d1e520e2e5f2f0e82c20d1f2f0e8e1eee6e820e2edf3f6e82c20e2e5fef2fa20f120eceef0ff20f1f2f0e5ebe0ece820ede020f5f0e0e1f0fbff20efebfaeafb20c8e3eef0e5e2fb"},
}

var golden = []gost3411_2012Test{
	{"8e945da209aa869f0455928529bcae4679e9873ab707b55315f56ceb98bef0a7362f715528356ee83cda5f2aac4c6ad2ba3a715c1bcd81cb8e9f90bf4c1c1a8a", ""},
	{"8b2a40ecab7b7496bc4cc0f773595452baf658849b495acc3ba017206810efb00420ccd73fb3297e0f7890941b84ac4a8bc27e3c95e1f97c094609e2136abb7e", "a"},
	{"2f2afbc8d92a31e2ecff5c490f2ca5db183d1e8a72fb1e755eaaa76d346b80096052bbe8ddab2f1639e5e50e4e8c1f8222e56155a1331147ca3638fad4e412f6", "ab"},
	{"28156e28317da7c98f4fe2bed6b542d0dab85bb224445fcedaf75d46e26d7eb8d5997f3e0915dd6b7f0aab08d9c8beb0d8c64bae2ab8b3c8c6bc53b3bf0db728", "abc"},
	{"41ab33a79b0a59f1ce5455ab74ce4249a6d9fd83e8c9f7a11f05107c172fac93e561e440e0355addc67abb779d805835433366902509725b2c1a3b7cab005790", "abcd"},
	{"c867aa7f3946ff1247ce937f49023871e400dd58e6615dc862597c018bb9c95200620b705624bd0f853521574d6a62721de7a433719b403b6173ad710f20b219", "abcde"},
	{"7dcfb784c3420a364072d501aff508aa6ee8b701ca6248f7a0cd9fca885734e2dad2111e4fadc0ecaa7b06558e1d3e04cd6022953cbc20a39a4ca6d8b98a85e9", "abcdef"},
	{"19d6b632fbb5e9c9d5094ce1ed11816fabb66ce6d7d3437b669942e2bbdf238aaa89c24341c1d17a3c1a2c14f77a2613a6ef27ce5f3326b3343fc26bbd884272", "abcdefg"},
	{"12479adc37d094e19e440caeac75ada8d18fc9557a37cb3d841f4bd75f49c10dca5c3834a145a9de09aa0180fe5e32cd0d06d38320f0607777954dab3c9fdd22", "abcdefgh"},
	{"cff3fab7cb56e53c2ff90eb85b5c6c077b34fc797f7f8559d735885561bfb15c578d72d03a1af2f0d239acc49f6918f0258abefa0b5f3b1d622f49aa9354531a", "abcdefghi"},
	{"8de12f9f1484929bdfc772601821dfc2d3594b25990031275dd285ac3d1d91f4306b50c124bc76e529f2ef0d2f9fc005ef57af1fce3577a6901e7e39f6142c7c", "abcdefghij"},
	{"a9b8eef57cc0b453508e09e69b458d9d352574952f9f3649c1f4384e0f3d3f2021e52338e19a0ac12c583354332879f65489bf649422447866cdec5394c54b4b", "Discard medicine more than two years old."},
	{"21f8b0bb865a83aa4529694fe902966573f09daec951cb95ef1bd3ad34326e6b29c51c171488c929954db466b56de15010996d31b58075006fd984605db225f8", "He who has a shady past knows that nice guys finish last."},
	{"52418d95fe3adc7b00af2a7d2c6b2978f2181194bd0e3f1d6a13c1b3de407c0c91558fde7460175844295ad87997c62a29cbbd76a3b2db714fdf580579e23b12", "I wouldn't marry him with a ten foot pole."},
	{"451a1ab82cdfcfae91b6a7f066916ebc935fbcd12d812548e731fc1a211f851d0000c97904d2a4958d76e97bd6814aa07968f90ff798af2e739d2200c6113c5c", "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{"afe25f356188a42a2656df4143f80282b9cdf2c13d10b1fb31329caa0b8c8a7f12c528c96645891504e53e13e524026088dec6d01e09faf8c332e16d3eb7092f", "The days of the digital watch are numbered.  -Tom Stoppard"},
	{"9645d8f327bf512e80db70eeca4b3be1a81e363f75afd88bd2cbc858fe61f991973cb25517e2508d642f61765af1b4e44e54069e2045f35ad1ffff9ef0940eb9", "Nepal premier won't resign."},
	{"a221ab7acc6bee331ce77c35b199dbdc78710b728e299b68c6e419a991355cea6145f114257c2995e33d05ac92aa9fb81b62bffefef7750faf29bf1c190aa901", "For every action there is an equal and opposite government program."},
	{"eaaae7d73f438de8c6f090bc5513885b0ef4b8da5f33a5500d946245a2493eb9ff0da0db0c55077e2929e3f1c77a63e1977e6b207d943f33a9d3af4ceedf5329", "His money is twice tainted: 'taint yours and 'taint mine."},
	{"ff7bede66ea3915d2c8e07b5396e3b178e6a0441c5840d67ca52b3dbf52841a6b9433092353f324c66de155fffd34da33c6f14e4b44ecd7372306659ab63c49e", "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{"e47b05f616f78d5d58e121692723507cc517855ca2b4661ae576f0b98e2eb6818dc5ac5633babd2eced7a369c29701714048647fff18042702a91a3fd1c48175", "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{"4fd43764ac4365d8b7c9c7653e56e948a1f5ddc117824854ebafd0e4b7d6291a5ad8f6f449d9edbe4c7ec7440f9be44fbdbf03fcc4f0fe0784846652069796a0", "size:  a.out:  bad magic"},
	{"d02729c9d3e8bc3f71c202bd9e5b3be6914217d2bc91c5f494a52cfd2715c9676b87390deb1de041b377f861519bd5549776f9530f7888d414bbc96346c7410a", "The major problem is with sendmail.  -Mark Horton"},
	{"20f44d8db9080767fc8d8db095dd80e0a33abc7be2da8f3b1675fc9e86ca13040d6adb1c8f5f28c490bd2ff44f412f13fe8e923d626e56372490310d893a9311", "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{"c37be7eac9c04628460b39c1950ed98184618f1d58d18088eb1d16762a392006588d22b7cbc27b812772d015d5d7adb78a3bffdecf324f13343dacad607f04d0", "If the enemy is within range, then so are you."},
	{"d86764908ddaac43d4a911c9b9826efbeded34b70627376d3a681ad109446917b06bdea30fd8d5ebcc8d0719d3aa4dc1aff8616b9a0898d0d1a02b58696deadc", "It's well we cannot hear the screams/That we create in others' dreams."},
	{"67095b50aa1a29f27e61e1de6d0d1cd4e1d45f74843fa6fb0a54390f79c22423ca1ff9c3eed99e390e68d6b92a33641f4c5602a2e62c538a11786ea4b6b9a1d2", "You remind me of a TV show, but that's all right: I watch it anyway."},
	{"8247f2c68db564cb7fcb4e205aa9e7f41cdddfa6dfd3c87e9d9d2473fb1f6b291f30d6c6cfb82609877d7a3cffa8d9babde1d8fa8e3a0707945e70c0ebdd6434", "C is as portable as Stonehedge!!"},
	{"5c8d0d08d060112790d55800c72cf4ea0c869c610b42f2995b2892d05bb8ef143e99831eae48a5a7f54b6fe33b3053fd4e7b760f4f615094bc4d36fda7149cd5", "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{"3786c225ccc1c77702bbf421c4c47719e8a4405b9cf607615f8d6846a1598a4cf22ef947ed0c262006b30eb98d3562462293de794acb62fbf4184c8edc226edb", "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{"1528108c602480eabe2d028550dd66f5cdd83ac43256b5666d2e11634828eec7d3e82900e58748bd2bb56f7146043117495cd636bafadc854fc8468f1c2e6227", "How can you write a big system without C++?  -Paul Glick"},
}

var golden256 = []gost3411_2012Test{
	{"3f539a213e97c802cc229d474c6aa32a825a360b2a933a949fd925208d9ce1bb", ""},
	{"ba31099b9cc84ec2a671e9313572378920a705b363b031a1cb4fc03e01ce8df3", "a"},
	{"0823c97ef0de722c918994944334e2a312889851a9667468b019b8582b35b590", "ab"},
	{"4e2919cf137ed41ec4fb6270c61826cc4fffb660341e0af3688cd0626d23b481", "abc"},
	{"8ec485bad9dcdeae7c7201d85cc6a9faa5f573d4bf2e1639bdd5b54c84a59e45", "abcd"},
	{"dda887af02d8c39e0138bd4b95f8cf0ddaf7cd4637fcb94d55bb4003339ec01e", "abcde"},
	{"f966c3d39587addd9f724ea4c0c6c7e6e044c02fbb3a10ca745658f5819dae0c", "abcdef"},
	{"5ad83263aee66b0620b8b6abe8432b6283360a3d60bc287b7e2bf41b19b318cf", "abcdefg"},
	{"1983d6d1da5171c1d9cb29fe9d6128f699a74a321924cddf724d6cc8326e1887", "abcdefgh"},
	{"44f796c0e845560b25ba92ceaa90e8b4b5daa888ee592d3bd9d07772b1d08e44", "abcdefghi"},
	{"e5952c9f46800d64580d6d43022ce67eddb227bb7a5b3f7a87ea3ff768d6ba9f", "abcdefghij"},
	{"86b0fa34d4c7ca23faea1decd157f10c8c0a91824f4ce6ce60679752b66902c9", "Discard medicine more than two years old."},
	{"659f2727409a801290b0f000173e5418ab9a47dcca7970b9c02d23e3c60f2b5a", "He who has a shady past knows that nice guys finish last."},
	{"0102ed219892003cb16e390c8520d32dde1df16ff447ded7d103f5f575ac2c5f", "I wouldn't marry him with a ten foot pole."},
	{"e9f1451455d2ea85bce8b52134d2c43a785647f5246f321008e312a5f4197dfd", "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{"55c880913db793be7de1528e8061464610f4860eb917442cd05b7af38c736669", "The days of the digital watch are numbered.  -Tom Stoppard"},
	{"855ca83c8a94a231ba7afa917afd509095153ba0f169ac710950282b40cfab30", "Nepal premier won't resign."},
	{"ca468b2f12917bf9f94b0b81db6ddd37c9101d419e1256b80a4ff79b7a95fe5c", "For every action there is an equal and opposite government program."},
	{"f91d24de8500ca3e5fc48a4dcc244071bfaee91e798e39dda9a6006b92a4d835", "His money is twice tainted: 'taint yours and 'taint mine."},
	{"ba160d6b72402f113a4226d76443c3f6d096d5a758ae588e1ca6ff8a22ace447", "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{"25e336ba51f11f372636cdfbb6c550a3b6c1729c3008017da93228de568c7624", "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{"6950d3a73836d7075d47cc7336b9eb231ef6338570f07a53e8e9adabe64e3514", "size:  a.out:  bad magic"},
	{"7e3424f7fc8007a7f680eb088cbd591dda9ab939d4826d6ec092acd702ade174", "The major problem is with sendmail.  -Mark Horton"},
	{"05e6cc07c91c554f7e921accbf52b691a988b65e7f00a9e949d811820b247a2f", "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{"9fec0cb244f0113eefe967e731d0e177eae712f35258d9a7e3d4322bb8d2ec45", "If the enemy is within range, then so are you."},
	{"76ff75e5d17237f86c0691ee811dd95e4cec3d3a3216a10bf616c3c14439bfc7", "It's well we cannot hear the screams/That we create in others' dreams."},
	{"3315e6e996ddc0be0325c6be815b03b180a953ffc773f3e94f36af19bffe08ac", "You remind me of a TV show, but that's all right: I watch it anyway."},
	{"3134945238ee2767b2a12f9a637f19d5cebaf9386592396222d41b7b3e5ea082", "C is as portable as Stonehedge!!"},
	{"541ab367d229e78f4ae2b4649e628269753d4e8b98b0370564c7391c8080dbda", "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{"971bcf1eb54bd9430ab2268d8bb9db3af9ac99064769fa515cd0ffd6d03d5183", "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{"bcc500b1f28310b5123f923c7ed0ef1dcbd218095a3ccf6bf1def24842574bbc", "How can you write a big system without C++?  -Paul Glick"},
}

func TestStandard(t *testing.T) {
	for i := 0; i < len(standard); i++ {
		g := standard[i]
		in, _ := hex.DecodeString(g.in)
		s := fmt.Sprintf("%x", Sum512(in))
		if s != g.out {
			t.Fatalf("Sum512 function: gost3411_2012_512(%s) = %s want %s", g.in, s, g.out)
		}
		c := New()
		for j := 0; j < 3; j++ {
			if j < 2 {
				c.Write(in)
			} else {
				c.Write(in[0 : len(in)/2])
				c.Sum(nil)
				c.Write(in[len(in)/2:])
			}
			s := fmt.Sprintf("%x", c.Sum(nil))
			if s != g.out {
				t.Fatalf("gost3411_2012_512[%d](%s) = %s want %s", j, g.in, s, g.out)
			}
			c.Reset()
		}
	}
	for i := 0; i < len(standard256); i++ {
		g := standard256[i]
		in, _ := hex.DecodeString(g.in)
		s := fmt.Sprintf("%x", Sum256(in))
		if s != g.out {
			t.Fatalf("Sum256 function: gost3411_2012_256(%s) = %s want %s", g.in, s, g.out)
		}
		c := New256()
		for j := 0; j < 3; j++ {
			if j < 2 {
				c.Write(in)
			} else {
				c.Write(in[0 : len(in)/2])
				c.Sum(nil)
				c.Write(in[len(in)/2:])
			}
			s := fmt.Sprintf("%x", c.Sum(nil))
			if s != g.out {
				t.Fatalf("gost3411_2012_256[%d](%s) = %s want %s", j, g.in, s, g.out)
			}
			c.Reset()
		}
	}
}

func TestGolden(t *testing.T) {
	for i := 0; i < len(golden); i++ {
		g := golden[i]
		s := fmt.Sprintf("%x", Sum512([]byte(g.in)))
		if s != g.out {
			t.Fatalf("Sum512 function: gost3411_2012_512(%s) = %s want %s", g.in, s, g.out)
		}
		c := New()
		for j := 0; j < 3; j++ {
			if j < 2 {
				io.WriteString(c, g.in)
			} else {
				io.WriteString(c, g.in[0:len(g.in)/2])
				c.Sum(nil)
				io.WriteString(c, g.in[len(g.in)/2:])
			}
			s := fmt.Sprintf("%x", c.Sum(nil))
			if s != g.out {
				t.Fatalf("gost3411_2012_512[%d](%s) = %s want %s", j, g.in, s, g.out)
			}
			c.Reset()
		}
	}
	for i := 0; i < len(golden256); i++ {
		g := golden256[i]
		s := fmt.Sprintf("%x", Sum256([]byte(g.in)))
		if s != g.out {
			t.Fatalf("Sum256 function: gost3411_2012_256(%s) = %s want %s", g.in, s, g.out)
		}
		c := New256()
		for j := 0; j < 3; j++ {
			if j < 2 {
				io.WriteString(c, g.in)
			} else {
				io.WriteString(c, g.in[0:len(g.in)/2])
				c.Sum(nil)
				io.WriteString(c, g.in[len(g.in)/2:])
			}
			s := fmt.Sprintf("%x", c.Sum(nil))
			if s != g.out {
				t.Fatalf("gost3411_2012_256[%d](%s) = %s want %s", j, g.in, s, g.out)
			}
			c.Reset()
		}
	}
}

func TestSize(t *testing.T) {
	c := New()
	if got := c.Size(); got != Size {
		t.Errorf("Size = %d; want %d", got, Size)
	}
	c = New256()
	if got := c.Size(); got != Size256 {
		t.Errorf("New256.Size = %d; want %d", got, Size256)
	}
}

func TestBlockSize(t *testing.T) {
	c := New()
	if got := c.BlockSize(); got != BlockSize {
		t.Errorf("BlockSize = %d want %d", got, BlockSize)
	}
}

var bench = New()
var buf = make([]byte, 8192)

func benchmarkSize(b *testing.B, size int) {
	b.SetBytes(int64(size))
	sum := make([]byte, bench.Size())
	for i := 0; i < b.N; i++ {
		bench.Reset()
		bench.Write(buf[:size])
		bench.Sum(sum[:0])
	}
}

func BenchmarkHash8Bytes(b *testing.B) {
	benchmarkSize(b, 8)
}

func BenchmarkHash1K(b *testing.B) {
	benchmarkSize(b, 1024)
}

func BenchmarkHash8K(b *testing.B) {
	benchmarkSize(b, 8192)
}
